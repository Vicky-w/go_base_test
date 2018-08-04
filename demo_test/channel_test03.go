package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//c:=make(chan bool)
	//for i :=0 ; i<10 ; i++{
	//	go Go(c,i)
	//}
	//<- c

	c := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go Go(c, i)
	}
	for i := 0; i < 10; i++ {
		v := <-c
		fmt.Println(v)
	}
}

//func Go(c chan bool, index int) {
//	a := 1
//	for i := 0; i < 10000000; i++ {
//		a += i
//	}
//	fmt.Println(index, a)
//	if index == 9 {   多CUP不能使用尾部index来判断
//		c <- true
//	}
//
//}
func Go(c chan bool, index int) { //只保证了输出数量 不保证顺序
	a := 1
	for i := 0; i < 10000000; i++ {
		a += i
	}
	fmt.Println(index, a)
	c <- true
}
