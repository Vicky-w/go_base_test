package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	path string //读取文件路径
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

func (r *ReadFromFile) Read(rc chan []byte) {
	//读取文件
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error : %s", err.Error()))
	}
	//从文件末尾开始逐行读取文件
	f.Seek(0, 2)
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF { //日志产生  结尾判断
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes errors : %s", err.Error()))
		}
		//rc <- line
		//去除换行符
		rc <- line[:len(line)-1]
	}
}
func (w *WriteToInfluxDB) Write(wc chan *Message) {
	//初始化influxdb client
	// 从Write Channel中读取监控数据
	//构造数据并写入influxdb
	/*
		Influxdb 是一个开源的时序型的数据库,使用Go语言编写，被广泛应用于存储系统的监控数据，IoT行业的实时数据场景。有一下特征：

		部署简单，无外部依赖
		内置http支持， 使用http读写
		类sql的灵活查询（max、min、sum 等）
	*/
	for v := range wc {
		fmt.Println(v)
	}
}

func (l *LogProcess) Process() {
	//从Read Channel 中读取每行日志数据
	//正则提取所需的监控数据（path status method 等）
	//写入Write Channel
	/**
	172.0.0.12 - -[04/Mar/2018:13:49:52 +0000] http "GET /vickywang?query=t HTTP/1.0" 200 2133 "-"
	"KeepAliveClient" "-" 1.005 1.854

	([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([~\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)
	*/
	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([~\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for v := range l.rc {
		ret := r.FindStringSubmatch(string(v))
		if len(ret) != 14 {
			log.Println("FindStringSubmatch fail:", string(v))
			continue
		}
		message := &Message{}
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc) //6 1 2 3 4 5
		if err != nil {
			log.Println("ParseInLocation fail:", ret[4])
		}
		message.TimeLocal = t
		byteSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = byteSent
		//GET /vickywang?query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			log.Println("strings.Split fail", ret[6])
			continue
		}
		message.Method = ret[0]
		u, err := url.Parse(ret[1])
		if err != nil {
			log.Println("url parse fail:", err)
		}
		message.Path = u.Path
		message.Scheme = ret[5]
		message.Status = ret[7]
		upstreamTime, _ := strconv.ParseFloat(ret[12], 64)
		requestTime, _ := strconv.ParseFloat(ret[13], 64)
		message.UpstreamTime = upstreamTime
		message.RequestTime = requestTime
		l.wc <- message
	}
}

func main() {
	r := &ReadFromFile{
		path: "./tmp/access.log",
	}
	w := &WriteToInfluxDB{
		influxDBDsn: "username&password",
	}
	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan *Message),
		read:  r,
		write: w,
	}
	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)
	time.Sleep(30 * time.Second)
}

/*
go run log_process.log

echo VickyWang01 >> access.log
echo VickyWang02 >> access.log
*/
