package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	_ "time"
)

const bufsize = 1024

var kvFile = path.Join(os.Getenv("HOME"), ".gokv.json")

var kv map[string]interface{}

func panicIfErr(e error) {
	if e != nil {
		panic(e)
	}
}

func handle(c io.ReadWriteCloser) {
	defer func() {
		//<-time.After(3* time.Second)
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

func loadKv() {
	data, err := ioutil.ReadFile(kvFile)
	panicIfErr(err)
	panicIfErr(json.Unmarshal(data, &kv))
}

func main() {
	loadKv()
	//defer persistKv()
	fmt.Println("starting server on localhost 4000")
	l, err := net.Listen("tcp", ":4000")
	panicIfErr(err)
	for {
		c, err := l.Accept()
		panicIfErr(err)
		go handle(c)
	}
}
