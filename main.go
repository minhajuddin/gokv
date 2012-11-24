package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

const bufsize = 1024

//file where the serialized json is persisted
var kvFile = path.Join(os.Getenv("HOME"), ".gokv.json")

//in memory map of key value store
var kv map[string]interface{}

//helper method
func panicIfErr(e error) {
	if e != nil {
		log.Panic(e)
	}
}

//core handler
//if the input is just one line it tries to get the
//value for the given key
//if it is more than one line, it assumes that the 
//first line is a key and the rest is the value
//it tries to store the data in this scenario
func handle(c io.ReadWriteCloser) {
	defer log.Println("connection closed")
	defer c.Close()
	buf := make([]byte, bufsize)
	nr, _ := io.ReadFull(c, buf)
	key := strings.Trim(string(buf[:nr]), "\n")
	idx := strings.Index(key, "\n")
	//if it has a newline, the first line is the key
	if idx > -1 {
		kv[key[:idx]] = strings.Trim(key[idx:], "\n")
		return
	}
	if v, ok := kv[key]; ok {
		fmt.Fprintln(c, v)
		return
	}
	log.Printf("key for '%s' not found\n", key)
	fmt.Fprintln(c, "<NULL>")
}

//loads the key value data from the persistence file
//when the server is started
func loadKv() {
	data, err := ioutil.ReadFile(kvFile)
	panicIfErr(err)
	panicIfErr(json.Unmarshal(data, &kv))
}

//persists the data to the persistence file when the
//server shuts down
func persistKv() {
	bytes, err := json.Marshal(kv)
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile(kvFile, bytes, 0600)
	log.Println("all key values persisted")
}

func handleSysSignals(l *net.Listener) {
	//signal handling
	//this code allows us to handle the SIGINT and SIGTERM signals
	//gracefully, we persist the file before we shutdown 
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//wait for a SIGINT or SIGTERM signal on a 
	//different thread
	go func() {
		log.Printf("received a %v\n", <-sigs)
		//make sure that the listener is closed before graceful exit
		(*l).Close()
	}()
}

func main() {
	//load persistence file 
	loadKv()
	//save persistence file
	defer persistKv()

	log.Println("starting server on localhost 4000")
	l, err := net.Listen("tcp", ":4000")
	panicIfErr(err)

	handleSysSignals(&l)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("error in accept", err)
			break
		}
		log.Println("accepted connection")
		go handle(c)
	}
}
