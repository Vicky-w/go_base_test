package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Servers_write struct {
	XMLName xml.Name       `xml:"servers"`
	Version string         `xml:"version,attr"`
	Svs     []server_write `xml:"server"`
}

type server_write struct {
	ServerName string `xml:"serverName"`
	ServerIP   string `xml:"serverIP"`
}

func main() {
	v := &Servers_write{Version: "1"}
	v.Svs = append(v.Svs, server_write{"Shanghai_VPN", "127.0.0.1"})
	v.Svs = append(v.Svs, server_write{"Beijing_VPN", "127.0.0.2"})
	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	os.Stdout.Write([]byte(xml.Header)) //xml.MarshalIndent或者xml.Marshal输出的信息都是不带XML头

	os.Stdout.Write(output)

	/*
		tag中含有"a>b>c"，那么就会循环输出三个元素a包含b，b包含c，例如如下代码就会输出
		FirstName string   `xml:"name>first"`
		LastName  string   `xml:"name>last"`

		<name>
		<first>Asta</first>
		<last>Xie</last>
		</name>
	*/
}
