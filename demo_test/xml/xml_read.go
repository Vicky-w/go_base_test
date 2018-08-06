package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

/**
如果某个struct字段的tag定义 型如"a>b>c",则解析的时候，会将xml结构a下面的b下面的c元素的值赋值给该字段。
如果某个struct字段的tag定义了"-",那么不会为该字段解析匹配任何xml数据。
如果struct字段后面的tag定义了",any"，如果他的子元素在不满足其他的规则的时候就会匹配到这个字段
*/
type Recurlyservers struct {
	XMLName xml.Name `xml:"servers"`
	Version string   `xml:"version,attr"` //",attr" 属性的值赋值给该字段
	Svs     []server `xml:"server"`
	/*
		如果struct的一个字段是string或者[]byte类型且它的tag含有",innerxml"，Unmarshal将会将此字段所对应的元素内所有内嵌的原始xml
	*/
	Description string `xml:",innerxml"`
}

type server struct {
	XMLName    xml.Name `xml:"server"`
	ServerName string   `xml:"serverName"`
	ServerIP   string   `xml:"serverIP"`
}

func main() {
	file, err := os.Open("test.xml") // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	v := Recurlyservers{} //目前支持struct、slice和string
	/*
		XML本质上是一种树形的数据格式，而我们可以定义与之匹配的go 语言的 struct类型，然后通过xml.Unmarshal来将xml中的数据解析成对应的struct对象。
	*/
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Println(v)
	fmt.Println(v.Description)
}
