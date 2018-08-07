package main

import (
	"fmt"
	"html/template"
	"os"
)

type Friend struct {
	Fname string
}

type Person struct {
	UserName string
	Emails   []string
	Friends  []*Friend
}

func main() {
	f1 := Friend{Fname: "Dodo Wu"}
	f2 := Friend{Fname: "JingWen Zhou"}
	t := template.New("fieldname example")
	t, _ = t.Parse(`hello {{.UserName}}!
			{{range .Emails}}
				an email {{.}}
			{{end}}
			{{with .Friends}}
			{{range .}}
				my friend name is {{.Fname}}
			{{end}}
			{{end}}
			`)

	fmt.Println("The first one parsed OK.")

	p := Person{UserName: "Vicky Wang",
		Emails:  []string{"585288507@qq.com", "585288507@qq.com"},
		Friends: []*Friend{&f1, &f2}}
	t.Execute(os.Stdout, p)

	tOk := template.New("first")
	//检查是否正确
	template.Must(tOk.Parse(" some static text /* and a comment */")) //panic: template: check parse error with Must:1: unexpected "}" in command
}
