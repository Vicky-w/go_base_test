package main

import (
	"fmt"
	"runtime"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		fmt.Println(s)
		runtime.Gosched() //表示让CPU把时间片让给别人,下次某个时候继续恢复执行该goroutine
	}
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU()) 告诉调度器同时使用多个线程
	go say("world") //开一个新的Goroutines执行
	say("hello")    //当前Goroutines执行
}
