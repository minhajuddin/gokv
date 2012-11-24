package main

import (
	"fmt"
	"net"
	"io"
	"os"
	"strings"
	"time"
)

const bufsize = 1024

//TODO: should make this persistent
var kv = map[string]interface{}{
	"name" : "Khaja Minhajuddin",
	"blog" : "http://minhajuddin.com",
	"cpa" : "goserve .",
}

func panicIfErr(e error){
	if e != nil {
		panic(e)
	}
}

func handle(c io.ReadWriteCloser) {
	defer func(){
		<-time.After(3* time.Second)
		println("closed")
		c.Close()
	}()
	buf := make([]byte, bufsize)
	//io.ReadFull(c, buf)
	nr, _ := io.ReadFull(c, buf)
	key := strings.Trim(string(buf[:nr]), "\n\r")
	if v, ok := kv[key]; ok {
		fmt.Fprintln(c, v)
		return
	}
	fmt.Fprintf(os.Stderr, "key for '%s' not found\n", key)
	fmt.Fprintln(c, "<NULL>")
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
