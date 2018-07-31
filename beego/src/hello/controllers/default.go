package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	//c.Ctx.WriteString("hello world")
	//beego.AppConfig.String("appname")
	//strconv.Itoa()
	beego.Informational("trace test1")
	beego.Notice("info test1")
	beego.SetLevel(beego.LevelInformational)
	beego.Informational("trace test2")
	beego.Notice("info test2")
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
