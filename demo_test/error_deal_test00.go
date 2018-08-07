package main

import (
	"errors"
	"fmt"
)

/**
error是一个内置的接口类型，我们可以在/builtin/包下面找到相应的定义。而我们在很多内部包里面用到的 error是errors包下面的实现的私有结构errorString
*/
func Sqrt(f float64) (float64, error) {
	//可以通过errors.New把一个字符串转化为errorString，以得到一个满足接口error的对象
	if f < 0 {
		return 0, errors.New("math: square root of negative number")
	} else {
		return 200, nil
	}
}

func main() {
	f, err := Sqrt(-1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f)
}
