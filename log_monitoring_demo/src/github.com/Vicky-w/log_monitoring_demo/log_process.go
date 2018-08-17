package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan *Message)
}

type LogProcess struct {
	rc     chan []byte
	wc     chan *Message
	reader Reader
	writer Writer
}

type WriteToInfluxDB struct {
	//influxDBDsn string // influx data source
	batch      uint16
	retry      uint8
	influxConf *InfluxConf
}

type Message struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}

// 系统状态监控
type SystemInfo struct {
	HandleLine   int       `json:"handleLine"`   // 总处理日志行数
	Tps          float64   `json:"tps"`          // 系统吞出量
	ReadChanLen  int       `json:"readChanLen"`  // read channel 长度
	WriteChanLen int       `json:"writeChanLen"` // write channel 长度
	RunTime      string    `json:"runTime"`      // 运行总时间
	ErrInfo      ErrorInfo `json:"errInfo"`      // 错误信息
	//ErrNum       int     `json:"errNum"`       // 错误数
}
type ReadFromTail struct {
	// 文件读取
	inode uint64
	fd    *os.File
	path  string // 读取文件的路径
}
type InfluxConf struct {
	//database 配置
	Addr, Username, Password, Database, Precision string
}
type Monitor struct {
	listenPort string
	startTime  time.Time
	tpsSli     []int
	systemInfo SystemInfo
}

type ErrorInfo struct {
	//错误信息
	ReadErr    int `json:"readErr"`
	ProcessErr int `json:"processErr"`
	WriteErr   int `json:"writeErr"`
}
type TypeMonitor int

const (
	TypeHandleLine TypeMonitor = iota
	TypeReadErr
	TypeProcessErr
	TypeWriteErr
)

var (
	path, influxDsn, listenPort string
	processNum, writeNum        int
	TypeMonitorChan             = make(chan TypeMonitor, 200)
)

//构造函数
func NewReader(path string) (Reader, error) {
	var stat syscall.Stat_t //TODO
	if err := syscall.Stat(path, &stat); err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &ReadFromTail{
		inode: stat.Ino,
		fd:    f,
		path:  path,
	}, nil
}

func NewWriter(influxDsn string) (Writer, error) {
	influxDsnSli := strings.Split(influxDsn, "@")
	if len(influxDsnSli) < 5 {
		return nil, errors.New("param influxDns err")
	}
	return &WriteToInfluxDB{
		batch: 50,
		retry: 3,
		influxConf: &InfluxConf{
			Addr:      influxDsnSli[0],
			Username:  influxDsnSli[1],
			Password:  influxDsnSli[2],
			Database:  influxDsnSli[3],
			Precision: influxDsnSli[4],
		},
	}, nil
}

func NewLogProcess(reader Reader, writer Writer) *LogProcess {
	return &LogProcess{
		rc:     make(chan []byte, 200),
		wc:     make(chan *Message, 200),
		reader: reader,
		writer: writer,
	}
}

func (m *Monitor) start(lp *LogProcess) {
	/*
		总处理日志行数
		系统吞吐量
		read channel 长度
		write channnel 长度
		运行总时间
		错误数
	*/
	go func() { //消费监控数据
		for n := range TypeMonitorChan {
			switch n {
			case TypeHandleLine:
				m.systemInfo.HandleLine += 1
			case TypeReadErr:
				m.systemInfo.ErrInfo.ReadErr += 1
			case TypeProcessErr:
				m.systemInfo.ErrInfo.ProcessErr += 1
			case TypeWriteErr:
				m.systemInfo.ErrInfo.WriteErr += 1
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 5) //定时器 tps
	go func() {
		for {
			<-ticker.C //每5秒触发一次
			m.tpsSli = append(m.tpsSli, m.systemInfo.HandleLine)
			if len(m.tpsSli) > 2 { //做切割 总数始终为2
				m.tpsSli = m.tpsSli[1:]
			}
		}
	}()

	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) { //暴露接口
		/*
			m.data.RunTime = time.Now().Sub(m.startTime).String() // To compute t-d for a duration d, use t.Add(-d).  当前时间减开始时间
				m.data.ReadChanLen = len(lp.rc)
				m.data.WriteChanLen = len(lp.wc)
				if len(m.tpsSli) >= 2 { //tps 计算
					m.data.Tps = float64(m.tpsSli[1]-m.tpsSli[0]) / 5 //定时器设置为5秒
				}
				ret, _ := json.MarshalIndent(m.data, "", "\t")
				io.WriteString(writer, string(ret))
		*/
		io.WriteString(writer, m.systemStatus(lp))
	})
	//http.ListenAndServe(":50081", nil)
	log.Fatal(http.ListenAndServe(":"+m.listenPort, nil))
}
func (m *Monitor) systemStatus(lp *LogProcess) string {
	d := time.Now().Sub(m.startTime)
	m.systemInfo.RunTime = d.String()
	m.systemInfo.ReadChanLen = len(lp.rc)
	m.systemInfo.WriteChanLen = len(lp.wc)
	if len(m.tpsSli) >= 2 {
		// return math.Trunc(float64(m.tpsSli[1]-m.tpsSli[0])/5*1e3+0.5) * 1e-3
		m.systemInfo.Tps = float64(m.tpsSli[1]-m.tpsSli[0]) / 5
	}
	res, _ := json.MarshalIndent(m.systemInfo, "", "\t")
	return string(res)
}
func (r *ReadFromTail) Read(rc chan []byte) {
	// 读取模块

	// 打开文件
	//f, err := os.Open(r.path)
	//if err != nil {
	//	panic(fmt.Sprintf("open file error:%s", err.Error()))
	//}

	// 从文件末尾开始逐行读取文件内容
	//f.Seek(0, 2)
	//rd := bufio.NewReader(f)

	defer close(rc)
	var stat syscall.Stat_t

	r.fd.Seek(0, 2) // seek 到末尾
	bf := bufio.NewReader(r.fd)
	for {
		line, err := bf.ReadBytes('\n')
		if err == io.EOF { //TODO 日志按天或按周做分割
			//文件名会变 如果文件名变了就去重新打开文件做操作
			//time.Sleep(500 * time.Millisecond)
			//continue
			if err := syscall.Stat(r.path, &stat); err != nil {
				// 文件切割，但新文件还没有生成
				time.Sleep(1 * time.Second)
			} else {
				nowInode := stat.Ino
				if nowInode == r.inode {
					// 无新的数据产生
					time.Sleep(1 * time.Second)
				} else {
					// 文件切割，重新打开文件
					r.fd.Close()
					fd, err := os.Open(r.path)
					if err != nil {
						panic(fmt.Sprintf("Open file err: %s", err.Error()))
					}
					r.fd = fd
					bf = bufio.NewReader(fd)
					r.inode = nowInode
				}
			}
			continue
		} else if err != nil {
			log.Printf("readFromTail ReadBytes err: %s", err.Error())
			TypeMonitorChan <- TypeReadErr
			continue
		}
		//rc <- line
		//去除换行符
		rc <- line[:len(line)-1]
	}
}

func (w *WriteToInfluxDB) Write(wc chan *Message) {
	// 写入模块

	//infSli := strings.Split(w.influxDBDsn, "@")

	//初始化influxdb client
	// 从Write Channel中读取监控数据
	//构造数据并写入influxdb
	/*
		Influxdb 是一个开源的时序型的数据库,使用Go语言编写，被广泛应用于存储系统的监控数据，IoT行业的实时数据场景。有一下特征：

		部署简单，无外部依赖
		内置http支持， 使用http读写
		类sql的灵活查询（max、min、sum 等）
	*/
	// https://github.com/influxdata/influxdb/tree/master/client
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     w.influxConf.Addr,
		Username: w.influxConf.Username,
		Password: w.influxConf.Password,
	})
	if err != nil {
		panic(fmt.Sprintf("influxdb NewHTTPClient err:%s", err.Error()))
	}

	//for v := range wc {
	for {
		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  w.influxConf.Database,
			Precision: w.influxConf.Precision,
		})
		if err != nil {
			panic(fmt.Sprintf("influxdb NewBatchPoints error:%s", err.Error()))
		}
		var count uint16

	Fetch:
		for v := range wc {
			// Create a point and add to batch
			// Tags: Path, Method, Scheme, Status
			// Fields: UpstreamTime, RequestTime, BytesSent
			tags := map[string]string{
				"Path":   v.Path,
				"Scheme": v.Scheme,
				"Status": v.Status,
			}

			fields := map[string]interface{}{
				"UpstreamTime": v.UpstreamTime,
				"RequestTime":  v.RequestTime,
				"BytesSent":    v.BytesSent,
			}

			pt, err := client.NewPoint("nginx_log", tags, fields, v.TimeLocal)
			if err != nil {
				TypeMonitorChan <- TypeWriteErr
				log.Println("influxdb NewPoint error:", err)
				continue
			}

			bp.AddPoint(pt)
			count++
			if count > w.batch {
				break Fetch
			}
		}

		var i uint8
		for i = 1; i <= w.retry; i++ {
			if err := c.Write(bp); err != nil {
				TypeMonitorChan <- TypeWriteErr
				log.Printf("influxdb Write error:%s, retry:%d", err.Error(), i)
				time.Sleep(1 * time.Second)
			} else {
				log.Println(w.batch, "point has written")
				break
			}
		}
		// Close client resources
		//if err := c.Close(); err != nil {
		//	log.Fatal(err)
		//}
		log.Println("write success!")
	}
}

func (l *LogProcess) Process() {
	// 解析模块

	//从Read Channel 中读取每行日志数据
	//正则提取所需的监控数据（path status method 等）
	//写入Write Channel
	/**
	172.0.0.12 - - [04/Mar/2018:13:49:52 +0000] http "GET /vickywang?query=t HTTP/1.0" 200 2133 "-" "KeepAliveClient" "-" 1.005 1.854
	*/

	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for v := range l.rc {
		TypeMonitorChan <- TypeHandleLine
		ret := r.FindStringSubmatch(string(v))
		//if len(ret) != 14 {
		if len(ret) < 13 { //???
			TypeMonitorChan <- TypeProcessErr //错误标记
			log.Println("FindStringSubmatch fail:", string(v))
			continue
		}

		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc)
		if err != nil {
			TypeMonitorChan <- TypeProcessErr //错误标记
			log.Println("ParseInLocation fail:", err.Error(), ret[4])
			continue
		}

		// GET /foo?query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			TypeMonitorChan <- TypeProcessErr //错误标记
			log.Println("strings.Split fail", ret[6])
			continue
		}

		u, err := url.Parse(reqSli[1])
		if err != nil {
			log.Println("url parse fail:", err)
			TypeMonitorChan <- TypeProcessErr //错误标记
			continue
		}
		//message.Method = reqSli[0]
		method := strings.TrimLeft(reqSli[0], "\"") //去除
		path := u.Path
		scheme := ret[5]
		status := ret[7]
		bytesSent, _ := strconv.Atoi(ret[8])
		upstreamTime, _ := strconv.ParseFloat(ret[12], 64)
		requestTime, _ := strconv.ParseFloat(ret[13], 64)

		//log.Println(message.TimeLocal)
		//log.Println(message.BytesSent)
		//log.Println(message.Method)
		//log.Println(message.Path)
		//log.Println(message.Scheme)
		//log.Println(message.Status)
		//log.Println(message.UpstreamTime)
		//log.Println(message.RequestTime)
		//l.wc <- message
		l.wc <- &Message{
			TimeLocal:    t,
			Path:         path,
			Method:       method,
			Scheme:       scheme,
			Status:       status,
			BytesSent:    bytesSent,
			UpstreamTime: upstreamTime,
			RequestTime:  requestTime,
		}
	}
}
func init() {
	flag.StringVar(&path, "path", "./tmp/access.log", "log file path")
	flag.StringVar(&influxDsn, "influxDsn", "http://172.16.17.137:8086@vickywang@vickywangpass@mydb@s", "influxDB dsn")
	flag.StringVar(&listenPort, "listenPort", "9193", "monitor port")
	flag.IntVar(&processNum, "processNum", 1, "process goroutine num")
	flag.IntVar(&writeNum, "writeNum", 1, "write goroutine num")
	flag.Parse()
}

func main() {
	reader, err := NewReader(path)
	if err != nil {
		panic(err)
	}

	writer, err := NewWriter(influxDsn)
	if err != nil {
		panic(err)
	}

	lp := NewLogProcess(reader, writer)

	go lp.reader.Read(lp.rc)
	for i := 0; i < 2; i++ { //并发执行
		go lp.Process()
	}

	for i := 0; i < 4; i++ { //并发执行
		go lp.writer.Write(lp.wc)
	}
	//监控程序
	m := &Monitor{
		listenPort: listenPort,
		startTime:  time.Now(),
	}
	go m.start(lp)
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)// syscall.SIGUSR1

	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("capture exit signal:", s)
			os.Exit(1)
			//case syscall.SIGUSR1: // 用户自定义信号
			//	log.Println(m.systemStatus(lp))
		default:
			log.Println("capture other signal:", s)
		}
	}
}

/*
go run log_process.log -path ./tmp/access.log -influxDsn http://172.16.17.137:8086@vickywang@vickywangpass@mydb@s

echo 172.0.0.12 - - [15/Aug/2018:20:34:07 +0000] https \"GET /bar HTTP/1.0\" 200 1261 \"-\" \"KeepAliveClient\" \"-\" - 0.017 >> access.log
echo VickyWang01 >> access.log
echo VickyWang02 >> access.log

监控需求：
  某个协议下的某个请求在某个请求方法的 QPS&响应时间&流量

Tags：Path, Method, Scheme, Status
Fields： UpstreamTime, RequestTime ，BytesSent
Time：TimeLocal

curl -XPOST "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE mydb"
*/
