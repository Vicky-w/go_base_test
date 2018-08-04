# Go 类型与变量

### 布尔型 `bool`

- 长度 1字节
- 取值范围  true false
- 注意事项  不可以数字代表true或false

### 整型 `int / uint `

- 根据运行平台可能为32位或者64位

### 8位整型   `int8 /uint8 `

- 长度1字节
- 取值范围  -128~127/ 0~255

### 字节型 `byte （uint8别名） `

### 16位整型  `int16 / uint16 `    
           
- 长度：2字节
- 取值范围：-2^16/2~2^16/2-1       0~2^16-1

### 32位整型 `int32 （rune）/uint32 `

- rune  unii  字符处理
- 长度：4字节
- 取值范围：-2^32/2~2^32/2-1       0~2^32-1

### 浮点型 `float32/float64 `

- 长度： 4 / 8 字节
- 小数位： 精确到 7 / 15 小数位

### 复数 `complex64/complex128 `

- 长度：8/16 字节

```
足够保存指针的32位或者64位整数型 ： uintptr
其他值类型：  array 、struct、string
引用类型： slice、map、chan
接口类型：inteface
函数类型：func

slice  是array 的 高层封装
类型零值      值类型位0 string为空字符串   引用为nil    bool 为false
```

```
查看类型 是否溢出
可以使用math包   的   例子    math.MaxInt32  math.MinInt32
```

```
var a int
a = 123      分支时使用
var b int = 321     确定类型
var c =321      不确定类型
d :=456   //全局变量不可使用      最简写法
```

- var ()  全局变量组

```
var (
aaa = “hello”
sss，bbb = 1，2
//  ccc := 3  不可省略 var
)
```

### 函数内部声明

```
var a,b,c,d int
a,b,c,d =1,2,3,4
var a,b,c,d int =1,2,3,4
var a,b,c,d =1,2,3,4
a,b,c,d :=1,2,3,4

var a float32 = 1.1
b :=int(a)
```
