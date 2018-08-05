package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	c := make(chan int)
	o := make(chan bool)
	go func() {
		for {
			select {
			case v := <-c:
				println(v)
			case <-time.After(3 * time.Second):
				println("timeout")
				o <- true
				break
			}
		}
	}()
	fmt.Println(runtime.NumGoroutine()) //返回正在执行和排队的任务总数
	<-o
}
