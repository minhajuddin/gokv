package main

import (
	"fmt"
	"net"
	"io"
	"time"
)

func panicIfErr(e error){
	if e != nil {
		panic(e)
	}
}

func handle(c io.ReadWriteCloser) {
	fmt.Fprintln(c, "This is awesome")
	<-time.After(time.Second * 5)
	c.Close()
}

func main(){
	fmt.Println("starting server on localhost 4000")
	l, err := net.Listen("tcp", ":4000")
	panicIfErr(err)
	for {
		c, err := l.Accept()
		panicIfErr(err)
		go handle(c)
	}
}
