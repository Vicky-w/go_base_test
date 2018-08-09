package main

type LogProcess struct {
	path        string //读取文件路径
	influxDBDsn string // influx data source
}

func (l *LogProcess) ReadFromFile() {
	//读取模块

}
func (l *LogProcess) Process() {
	//解析模块
}
func (l *LogProcess) WriteToInfluxDB() {
	//写入模块
}
func main() {

}
