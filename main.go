package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
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
	//if it has a newline, the first line is the key
	nr, _ := io.ReadFull(c, buf)
	key := strings.Trim(string(buf[:nr]), "\n")
	idx := strings.Index(key, "\n")
	if idx > -1 {
		kv[key[:idx]] = strings.Trim(key[idx:], "\n")
		return
	}
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

//TODO: add TERM signal handler here
func persistKv() {
	bytes, err := json.Marshal(kv)
	if err != nil {
		println(err)
	}
	ioutil.WriteFile(kvFile, bytes, 0600)
	println("persist complete")
}

func main() {
	loadKv()
	defer persistKv()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("starting server on localhost 4000")
	l, err := net.Listen("tcp", ":4000")
	panicIfErr(err)
	defer l.Close()
	go func() {
		fmt.Fprintf(os.Stderr, "received a %v\n", <-sigs)
		l.Close()
	}()
	for {
		c, err := l.Accept()
		if err != nil {
			println("errored in accept", err)
			break
		}
		println("accepted conn")
		go handle(c)
	}
}
