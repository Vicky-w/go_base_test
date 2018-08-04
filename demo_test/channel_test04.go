package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := sync.WaitGroup{} //值类型
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go Go2(&wg, i)
	}
	wg.Wait()

}
func Go2(wg *sync.WaitGroup, index int) {
	a := 1
	for i := 0; i < 10000000; i++ {
		a += i
	}
	fmt.Println(index, a)
	wg.Done()
}
