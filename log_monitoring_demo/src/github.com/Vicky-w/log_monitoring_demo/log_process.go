package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	TypeHandleLine = iota
	TypeErrNum
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan *Message)
}

type LogProcess struct {
	rc    chan []byte
	wc    chan *Message
	read  Reader
	write Writer
}

type ReadFromFile struct {
	path string // 读取文件的路径
}

type WriteToInfluxDB struct {
	influxDBDsn string // influx data source
}

type Message struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}

// 系统状态监控
type SystemInfo struct {
	HandleLine   int     `json:"handleLine"`   // 总处理日志行数
	Tps          float64 `json:"tps"`          // 系统吞出量
	ReadChanLen  int     `json:"readChanLen"`  // read channel 长度
	WriteChanLen int     `json:"writeChanLen"` // write channel 长度
	RunTime      string  `json:"runTime"`      // 运行总时间
	ErrNum       int     `json:"errNum"`       // 错误数
}
type ReadFromTail struct {
	// 文件读取
	inode uint64
	fd    *os.File
	path  string
}
type InfluxConf struct {
	//database 配置
	Addr, Username, Password, Database, Precision string
}
type Monitor struct {
	startTime time.Time
	data      SystemInfo
	tpsSli    []int
}
type ErrorInfo struct {
	//错误信息
	ReadErr    int `json:"readErr"`
	ProcessErr int `json:"processErr"`
	WriteErr   int `json:"writeErr"`
}

var TypeMonitorChan = make(chan int, 200)

func (m *Monitor) start(lp *LogProcess) {
	/*
		总处理日志行数
		系统吞吐量
		read channel 长度
		write channnel 长度
		运行总时间
		错误数
	*/
	go func() { //消费数据
		for n := range TypeMonitorChan {
			switch n {
			case TypeErrNum:
				m.data.ErrNum += 1
			case TypeHandleLine:
				m.data.HandleLine += 1
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 5) //定时器 tps
	go func() {
		for {
			<-ticker.C //每5秒触发一次
			m.tpsSli = append(m.tpsSli, m.data.HandleLine)
			if len(m.tpsSli) > 2 { //做切割 总数始终为2
				m.tpsSli = m.tpsSli[1:]
			}
		}
	}()

	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) { //暴露接口
		m.data.RunTime = time.Now().Sub(m.startTime).String() // To compute t-d for a duration d, use t.Add(-d).  当前时间减开始时间
		m.data.ReadChanLen = len(lp.rc)
		m.data.WriteChanLen = len(lp.wc)

		if len(m.tpsSli) >= 2 { //tps 计算
			m.data.Tps = float64(m.tpsSli[1]-m.tpsSli[0]) / 5 //定时器设置为5秒
		}

		ret, _ := json.MarshalIndent(m.data, "", "\t")
		io.WriteString(writer, string(ret))
	})

	http.ListenAndServe(":50081", nil)
}

func (r *ReadFromFile) Read(rc chan []byte) {
	// 读取模块

	// 打开文件
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error:%s", err.Error()))
	}

	// 从文件末尾开始逐行读取文件内容
	f.Seek(0, 2)
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF { //TODO 日志按天或按周做分割
			//文件名会变 如果文件名变了就去重新打开文件做操作
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes error:%s", err.Error()))
		}
		TypeMonitorChan <- TypeHandleLine //标记  总日志数量
		//rc <- line
		//去除换行符
		rc <- line[:len(line)-1]
	}
}

func (w *WriteToInfluxDB) Write(wc chan *Message) {
	// 写入模块

	infSli := strings.Split(w.influxDBDsn, "@")

	//初始化influxdb client
	// 从Write Channel中读取监控数据
	//构造数据并写入influxdb
	/*
		Influxdb 是一个开源的时序型的数据库,使用Go语言编写，被广泛应用于存储系统的监控数据，IoT行业的实时数据场景。有一下特征：

		部署简单，无外部依赖
		内置http支持， 使用http读写
		类sql的灵活查询（max、min、sum 等）
	*/
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     infSli[0],
		Username: infSli[1],
		Password: infSli[2],
	})
	if err != nil {
		log.Fatal(err)
	}

	for v := range wc {
		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  infSli[3],
			Precision: infSli[4],
		})
		if err != nil {
			log.Fatal(err)
		}

		// Create a point and add to batch
		// Tags: Path, Method, Scheme, Status
		tags := map[string]string{"Path": v.Path, "Method": v.Method, "Scheme": v.Scheme, "Status": v.Status}
		// Fields: UpstreamTime, RequestTime, BytesSent
		fields := map[string]interface{}{
			"UpstreamTime": v.UpstreamTime,
			"RequestTime":  v.RequestTime,
			"BytesSent":    v.BytesSent,
		}

		pt, err := client.NewPoint("nginx_log", tags, fields, v.TimeLocal)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)

		// Write the batch
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
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
		ret := r.FindStringSubmatch(string(v))
		if len(ret) != 14 {
			TypeMonitorChan <- TypeErrNum //错误标记
			log.Println("FindStringSubmatch fail:", string(v))
			continue
		}

		message := &Message{}
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc)
		if err != nil {
			TypeMonitorChan <- TypeErrNum //错误标记
			log.Println("ParseInLocation fail:", err.Error(), ret[4])
			continue
		}
		message.TimeLocal = t

		byteSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = byteSent

		// GET /foo?query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			TypeMonitorChan <- TypeErrNum //错误标记
			log.Println("strings.Split fail", ret[6])
			continue
		}
		message.Method = reqSli[0]

		u, err := url.Parse(reqSli[1])
		if err != nil {
			log.Println("url parse fail:", err)
			TypeMonitorChan <- TypeErrNum //错误标记
			continue
		}
		message.Path = u.Path

		message.Scheme = ret[5]
		message.Status = ret[7]

		upstreamTime, _ := strconv.ParseFloat(ret[12], 64)
		requestTime, _ := strconv.ParseFloat(ret[13], 64)
		message.UpstreamTime = upstreamTime
		message.RequestTime = requestTime
		//log.Println(message.TimeLocal)
		//log.Println(message.BytesSent)
		//log.Println(message.Method)
		//log.Println(message.Path)
		//log.Println(message.Scheme)
		//log.Println(message.Status)
		//log.Println(message.UpstreamTime)
		//log.Println(message.RequestTime)
		l.wc <- message
	}
}

func main() {
	var path, influxDsn string
	flag.StringVar(&path, "path", "./tmp/access.log", "read file path")
	flag.StringVar(&influxDsn, "influxDsn", "http://172.16.17.137:8086@vickywang@vickywangpass@mydb@s", "influx data source")
	flag.Parse()

	r := &ReadFromFile{
		path: path,
	}

	w := &WriteToInfluxDB{
		influxDBDsn: influxDsn,
	}

	lp := &LogProcess{
		rc:    make(chan []byte, 200), //配置缓存
		wc:    make(chan *Message, 200),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	for i := 0; i < 2; i++ { //并发执行
		go lp.Process()
	}

	for i := 0; i < 4; i++ { //并发执行
		go lp.write.Write(lp.wc)
	}
	//监控程序
	m := &Monitor{
		startTime: time.Now(),
		data:      SystemInfo{},
	}
	m.start(lp)
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
