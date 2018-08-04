package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	c := make(chan string)
	go func() {
		c <- "hello"
	}()

	go func() {
		word := <-c + " world"
		fmt.Println(word)
	}()
	time.Sleep(1 * time.Second)

	var n sync.WaitGroup
	for i := 0; i < 20; i++ {
		n.Add(1)
		go func(i int, n *sync.WaitGroup) {
			defer n.Done()
			time.Sleep(1 * time.Second)
			fmt.Printf("goroutine %d is running\n", i)
		}(i, &n)
	}
	n.Wait()
}
