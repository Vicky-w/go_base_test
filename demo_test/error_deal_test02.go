package main

func main() {
	/*
		package net
	    type Error interface {
	    error
	    Timeout() bool   // Is the error a timeout?
	    Temporary() bool // Is the error temporary?
	}
	*/
	/*
		在调用的地方，通过类型断言err是不是net.Error,来细化错误的处理，例如下面的例子，如果一个网络发生临时性错误，那么将会sleep 1秒之后重试
	*/
	//if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
	//	time.Sleep(1e9)
	//	continue
	//}
	//if err != nil {
	//	log.Fatal(err)
	//}
}
