package main

import (
	"fmt"
	"os"
)

func main() {
	var buffer [512]byte

	n, err := os.Stdin.Read(buffer[:])
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	fmt.Println("count:", n, ", msg:", string(buffer[:]))
}
