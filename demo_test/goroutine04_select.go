package main

import "fmt"

//select {
//case i := <-c:
// use i
//default:
// 当c阻塞的时候执行这里
//}

func fibonacci2(c, quit chan int) {
	x, y := 1, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
			fmt.Println("fibonacci2     x =", x, "     y=", y)
			c <- y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
func main() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("=============", i, "================")
			fmt.Println("read chan  x  ", <-c)
			fmt.Println("read chan  y  ", <-c)
		}
		quit <- 0
	}()
	fibonacci2(c, quit)
}
