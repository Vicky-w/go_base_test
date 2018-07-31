package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
)

const (
	_MYSQL_DRIVER = "mysql"
)

type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}
type Topic struct {
	Id               int64
	Uid              int64
	Title            string
	Content          string    `orm:"size(5000)"`
	Attachment       string
	Created          time.Time `orm:"index"`
	Updated          time.Time `orm:"index"`
	Views            int64     `orm:"index"`
	Author           string
	ReplyTime        time.Time `orm:"index"`
	ReplyCount       int64
	RepleyLastUserId int64
}

func RegisterDB() {
	//if !com.IsExist(_DB_NAME) {
	//	os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
	//	os.Create(_DB_NAME)
	//}
	//1.驱动类型
	orm.RegisterDriver(_MYSQL_DRIVER, orm.DRMySQL)
	//2.数据库配置
	dbHost := beego.AppConfig.String("db.host")
	dbPort := beego.AppConfig.String("db.port")
	dbUserName := beego.AppConfig.String("db.username")
	dbPwd := beego.AppConfig.String("db.pwd")
	dbDataBase := beego.AppConfig.String("db.database")
	//3. 数据库连接
	conn := dbUserName + ":" + dbPwd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbDataBase + "?charset=utf8"
	//4.注册默认数据库 30连接数
	orm.RegisterDataBase("default", _MYSQL_DRIVER, conn, 30, 30) //注册默认数据库
	//5.注册实体
	orm.RegisterModel(new(Category), new(Topic))
	//6.自动同步表结构
	orm.RunSyncdb("default", false, true)
}
