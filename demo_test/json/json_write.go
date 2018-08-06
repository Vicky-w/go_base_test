package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ServerWrite struct {
	ServerName string
	ServerIP   string
}

type ServersliceWrite struct {
	Servers []ServerWrite
}
type ServerWrite2 struct {
	ServerName string `json:"serverName"`
	ServerIP   string `json:"serverIP"`
}

type ServersliceWrite2 struct {
	Servers []ServerWrite2 `json:"servers"`
}
type ServerWrite3 struct {
	// ID 不会导出到JSON中
	ID int `json:"-"`

	// ServerName2 的值会进行二次JSON编码
	ServerName  string `json:"serverName"`
	ServerName2 string `json:"serverName2,string"`

	// 如果 ServerIP 为空，则不输出到JSON串中
	ServerIP string `json:"serverIP,omitempty"`
}

func main() {
	var s ServersliceWrite
	s.Servers = append(s.Servers, ServerWrite{ServerName: "Shanghai_VPN", ServerIP: "127.0.0.1"})
	s.Servers = append(s.Servers, ServerWrite{ServerName: "Beijing_VPN", ServerIP: "127.0.0.2"})
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Println(string(b))
	/**
	想输出小写不能将结构体的字段改成小写 必须使用 struct tag

	字段的tag是"-"，那么这个字段不会输出到JSON
	tag中带有自定义名称，那么这个自定义名称会出现在JSON的字段名中，例如上面例子中serverName
	tag中如果带有"omitempty"选项，那么如果该字段值为空，就不会输出到JSON串中
	如果字段类型是bool, string, int, int64等，而tag中带有",string"选项，那么这个字段在输出到JSON的时候会把该字段对应的值转换成JSON字符串
	*/
	s2 := ServerWrite3{
		ID:          3,
		ServerName:  `Go "1.0" `,
		ServerName2: `Go "1.0" `,
		ServerIP:    ``,
	}
	//JSON对象只支持string作为key
	b2, _ := json.Marshal(s2) //返回[]byte
	os.Stdout.Write(b2)
}
