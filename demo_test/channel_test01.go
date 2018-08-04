package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan []byte)

	go func(out chan<- []byte) {
		out <- []byte("VickyWang")
	}(ch)

	go func(in <-chan []byte) {
		fmt.Println(in)
	}(ch)

	time.Sleep(2 * time.Second)
}
