# goconfig 说明
goconfig是一个由Go语言开发的针对Windows下常见的INI格式的配置文件解析器。
该解析器在涵盖了所有INI文件操作的基础上，又针对Go语言实际开发过程中遇到的一些需求进行了扩展。

相对于其他INI文件解析器而言，该解析器最大的优势在于 ` 对注释的极佳支持 `  ；除此之外，
支持` 多个配置文件覆盖加载 ` 也是非常特别但好用的功能。

### 主要特性

- 提供与Windows API 一模一样的操作方式
- 支持`递归读取`分区
- 支持`自增键名`
- 支持对`注解的读与写`操作
- 支持`直接返回指定类型`的键值
- 支持`多个文件覆盖加载`

### 下载安装

- 通过 gopm 安装 `gopm get github.com/Unknwon/goconfig`
- 通过 go get 安装 `go get github.com/Unknwon/goconfig`

### 基本使用方法

- 加载配置文件 `cfg, err := goconfig.LoadConfigFile("config.ini")`
- 基本读写操作 
```
value, err := cfg.GetValue(goconfig.DEFAULT_SECTION, "key_default")
isInsert := cfg.SetValue(goconfig.DEFAULT_SECTION, "key_default", "这是新的值")
```
- 注释读写操作 
```
comment := cfg.GetSectionComments("super")
comment = cfg.GetKeyComments("super", "key_super")
v := cfg.SetSectionComments("super", "# 这是新的分区注释")
```
- 类型转换读取 `vInt, err := cfg.Int("must", "int")`
- Must系列方法 `vBool := cfg.MustBool("must", "bool")` 
- 删除制定键值 `ok := cfg.DeleteKey("must", "string")`
- 保存配置文件 `err = goconfig.SaveConfigFile(cfg, "conf_save.ini")`
- 多文件覆盖加载
```
cfg, err := goconfig.LoadConfigFile("conf.ini", "conf2.ini")
err = cfg.AppendFiles("conf3.ini")
```
- 配置文件重载 `err = cfg.Reload()`
- 为Must系列方法设置缺省值 `	vBool := cfg.MustBool("must", "bool404", true)`
- 递归读取键值 
- 子孙分区覆盖读取
- 自增键名获取
- 获取整个分区 `sec, err := cfg2.GetSection("auto increment")`
