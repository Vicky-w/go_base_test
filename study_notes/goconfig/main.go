package main

import (
	"github.com/gpmgo/gopm/modules/goconfig"
	"log"
)

func main() {
	// 创建并获取一个 ConfigFile 对象，以进行后续操作
	// 文件名支持相对和绝对路径
	cfg, err := goconfig.LoadConfigFile("conf.ini")
	if err != nil {
		log.Fatalf("无法加载配置文件：%s", err)
	}

	// 加载完成后所有数据均已存入内存，任何对文件的修改操作都不会影响到已经获取到的对象

	// >>>>>>>>>>>>>>> 基本读写操作 >>>>>>>>>>>>>>>

	// 对默认分区进行普通读取操作
	value, err := cfg.GetValue(goconfig.DEFAULT_SECTION, "key_default")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "key_default", err)
	}
	log.Printf("%s > %s: %s", goconfig.DEFAULT_SECTION, "key_default", value)

	// 对已有的键进行值重写操作，返回值为 bool 类型，表示是否为插入操作
	isInsert := cfg.SetValue(goconfig.DEFAULT_SECTION, "key_default", "这是新的值")
	log.Printf("设置键值 %s 为插入操作：%v", "key_default", isInsert)

	// 对不存在的键进行插入操作
	isInsert = cfg.SetValue(goconfig.DEFAULT_SECTION, "key_new", "这是新插入的键")
	log.Printf("设置键值 %s 为插入操作：%v", "key_new", isInsert)

	// 传入空白字符串也可直接操作默认分区
	value, err = cfg.GetValue("", "key_default")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "key_default", err)
	}
	log.Printf("%s > %s: %s", goconfig.DEFAULT_SECTION, "key_default", value)

	// 获取冒号为分隔符的键值
	value, err = cfg.GetValue("super", "key_super2")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "key_super2", err)
	}
	log.Printf("%s > %s: %s", "super", "key_super2", value)

	// <<<<<<<<<<<<<<< 基本读写操作 <<<<<<<<<<<<<<<

	// >>>>>>>>>>>>>>> 对注释进行读写操作 >>>>>>>>>>>>>>>

	// 获取某个分区的注释
	comment := cfg.GetSectionComments("super")
	log.Printf("分区 %s 的注释：%s", "super", comment)

	// 获取某个键的注释
	comment = cfg.GetKeyComments("super", "key_super")
	log.Printf("键 %s 的注释：%s", "key_super", comment)

	// 设置某个键的注释，返回值为 true 时表示注释被插入或删除（空字符串），false 表示注释被重写
	v := cfg.SetKeyComments("super", "key_super", "# 这是新的键注释")
	log.Printf("键 %s 的注释被插入或删除：%v", "key_super", v)

	// 设置某个分区的注释，返回值效果同上
	v = cfg.SetSectionComments("super", "# 这是新的分区注释")
	log.Printf("分区 %s 的注释被插入或删除：%v", "super", v)

	// <<<<<<<<<<<<<<< 对注释进行读写操作 <<<<<<<<<<<<<<<

	// >>>>>>>>>>>>>>> 自动转换和 Must 系列方法 >>>>>>>>>>>>>>>

	// 自动转换类型读取操作，直接返回指定类型，error 类型用于指示是否发生错误
	vInt, err := cfg.Int("must", "int")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "int", err)
	}
	log.Printf("%s > %s: %v", "must", "int", vInt)

	// Must 系列方法，一定返回某个类型的值；如果失败则返回零值
	vBool := cfg.MustBool("must", "bool")
	log.Printf("%s > %s: %v", "must", "bool", vBool)

	// 若键不存在则返回零值，此例应返回 false
	vBool = cfg.MustBool("must", "bool404")
	log.Printf("%s > %s: %v", "must", "bool404", vBool)

	// <<<<<<<<<<<<<<< 自动转换和 Must 系列方法 <<<<<<<<<<<<<<<

	// 删除键值，返回值用于表示是否删除成功
	ok := cfg.DeleteKey("must", "string")
	log.Printf("删除键值 %s 是否成功：%v", "string", ok)

	// 保存 ConfigFile 对象到文件系统，保存后的键顺序与读取时的一样
	err = goconfig.SaveConfigFile(cfg, "conf_save.ini")
	if err != nil {
		log.Fatalf("无法保存配置文件：%s", err)
	}

	// 创建并获取一个 ConfigFile 对象，以进行后续操作
	// 文件名支持相对和绝对路径，可指定多个文件名进行覆盖加载
	cfg2, err := goconfig.LoadConfigFile("conf.ini", "conf2.ini")
	if err != nil {
		log.Fatalf("无法加载配置文件：%s", err)
	}

	// 加载完成后所有数据均已存入内存，任何对文件的修改操作都不会影响到已经获取到的对象

	// 对默认分区进行普通读取操作
	value2, err := cfg2.GetValue(goconfig.DEFAULT_SECTION, "key_default")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "key_default", err)
	}
	log.Printf("%s > %s: %s", goconfig.DEFAULT_SECTION, "key_default", value2)

	// 若外部文件发生修改，可通过调用方法进行快速重载
	err = cfg2.Reload()
	if err != nil {
		log.Fatalf("无法重载配置文件：%s", err)
	}

	// 若在调用 Must 系列方法时发生错误，则可设置缺省值
	vBool2 := cfg2.MustBool("must", "bool404", true)
	log.Printf("%s > %s: %v", "must", "bool404", vBool2)

	// 可在操作中途追加配置文件
	err = cfg2.AppendFiles("conf3.ini")
	if err != nil {
		log.Fatalf("无法追加配置文件：%s", err)
	}

	// 进行递归读取键值
	value2, err = cfg2.GetValue("", "search")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "search", err)
	}
	log.Printf("%s > %s: %s", goconfig.DEFAULT_SECTION, "search", value2)

	// >>>>>>>>>>>>>>> 子孙分区覆盖读取 >>>>>>>>>>>>>>>

	// 以半角符号 . 为分隔符来表示多级别分区

	// 当子孙分区某个键存在时，会直接获取
	value2, err = cfg2.GetValue("parent.child", "age")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "age", err)
	}
	log.Printf("%s > %s: %s", "parent.child", "age", value2)

	// 当子孙分区某个键不存在时，会向父区按级寻找
	value2, err = cfg2.GetValue("parent.child", "sex")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "sex", err)
	}
	log.Printf("%s > %s: %s", "parent.child", "sex", value2)

	// <<<<<<<<<<<<<<< 子孙分区覆盖读取 <<<<<<<<<<<<<<<

	// 进行自增键名获取，凡是键名为半角符号 - 的在加载时均会被处理为自增
	// 自增范围限制在相同分区内
	// 为了方便展示，此处直接结合获取整个分区的功能并打印
	sec, err := cfg2.GetSection("auto increment")
	if err != nil {
		log.Fatalf("无法获取分区：%s", err)
	}
	log.Printf("%s : %v", "auto increment", sec)
	log.Print(sec["#1"])
}
