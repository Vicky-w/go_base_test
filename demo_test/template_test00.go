package main

import (
	"html/template"
	"os"
)

type Person struct {
	UserName string
	email    string //未导出的字段，首字母是小写的
}

func main() {
	t := template.New("fieldname example")
	//t, _ = t.Parse("hello {{.UserName}}!")
	t, _ = t.Parse("hello {{.UserName}}! {{.email}}")
	p := Person{UserName: "VickyWang", email: "595288507@qq.com"} //email 将不会输出
	t.Execute(os.Stdout, p)
}
