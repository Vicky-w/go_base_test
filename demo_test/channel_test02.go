package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	shutdown := make(chan struct{})
	var n sync.WaitGroup
	n.Add(1)
	go Running(shutdown, &n)
	n.Add(1)
	go ListenStop(shutdown, &n)
	n.Wait()
}

func Running(shutdown <-chan struct{}, n *sync.WaitGroup) {
	defer n.Done()
	for {
		select {
		case <-shutdown:
			// 一旦关闭channel，则可以接收到nil。
			fmt.Println("shutdown goroutine")
			return
		default:
			fmt.Println("I am running")
			time.Sleep(1 * time.Second)
		}
	}
}

func ListenStop(shutdown chan<- struct{}, n *sync.WaitGroup) {
	defer n.Done()
	os.Stdin.Read(make([]byte, 1))
	// 如果用户输入了回车则退出关闭channel
	close(shutdown)
}
