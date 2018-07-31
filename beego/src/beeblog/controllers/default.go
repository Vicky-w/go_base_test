package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"

	c.Data["TrueCond"] = true
	c.Data["FalseCond"] = false

	type u struct {
		Name string
		Age  int8
		Sex  int8
	}
	user := &u{
		Name: "VickyWang",
		Age:  25,
		Sex:  1,
	}
	c.Data["User"] = user
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	c.Data["Nums"] = nums
	c.Data["TplVar"]="hey guys"
	c.Data["Html"]="<div>Hello beego</div>"
	c.Data["Pipe"]="<div>Hello beego</div>"
}
