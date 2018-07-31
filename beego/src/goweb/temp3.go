package main

import (
	"net/http"
	"time"
	"io"
	"log"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	server := http.Server{
		Addr:        ":8080",
		Handler:     &myHandler2{},
		ReadTimeout: 5 * time.Second,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/hello"] = sayHello
	mux["/bye"] = sayBye
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

type myHandler2 struct {
}

func (*myHandler2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, "URL: "+r.URL.String())
}
func sayHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world, this is version 3")
}

func sayBye(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Bye bye, this is version 3")
}
