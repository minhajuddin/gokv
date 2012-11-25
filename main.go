//gokv is a basic key value stored which is persisted to disk in a json serialized form
//it supports 
//	+ storing a value for a key
//	+ retreiving a value for a key
//	- deleting a key value pair
//	- retreiving keys with a given prefix

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
	"runtime"
	"strings"
	"sync"
	"syscall"
)

const bufsize = 1024

var (
	//file where the serialized json is persisted
	kvFile = path.Join(os.Getenv("HOME"), ".gokv.json")

	//mutex used to lock/unlock access to the kv store
	mutex = &sync.Mutex{}
	//in-memory key value store
	kv map[string]interface{}
)

//helper method
//TODO: remove this, and handle errors where they occur
func panicIfErr(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

//core handler
//if the input is just one line it tries to get the
//value for the given key
//if it is more than one line, it assumes that the 
//first line is a key and the rest is the value
//it tries to store the data in this scenario
func handle(c net.Conn) {
	defer log.Println("connection closed for ", c.RemoteAddr())
	defer c.Close()
	buf := make([]byte, bufsize)
	nr, _ := io.ReadFull(c, buf)
	key := strings.Trim(string(buf[:nr]), "\n")
	idx := strings.Index(key, "\n")
	//if it has a newline, the first line is the key
	if idx > -1 {
		setValue(key[:idx], strings.Trim(key[idx:], "\n"))
		return
	}
	if v, ok := getValue(key); ok {
		fmt.Fprintln(c, v)
		return
	}
	log.Printf("key for '%s' not found for '%v'\n", key, c.RemoteAddr())
	fmt.Fprintln(c, "<NULL>")
}

//loads the key value data from the persistence file
//when the server is started
func loadKv() error {
	data, err := ioutil.ReadFile(kvFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &kv)
}

//read writes are done after locking the current go routine
func getValue(key string) (value interface{}, ok bool) {
	mutex.Lock()
	defer mutex.Unlock()
	value, ok = kv[key]
	return
}

func setValue(key string, value interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	kv[key] = value
	runtime.Gosched()
}

//persists the data to the persistence file when the
//server shuts down
func persistKv() error {
	mutex.Lock()
	bytes, err := json.Marshal(kv)
	mutex.Unlock()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(kvFile, bytes, 0600)
}

func handleSysSignals(l net.Listener) {
	//signal handling
	//this code allows us to handle the SIGINT and SIGTERM signals
	//gracefully, we persist the file before we shutdown 
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//wait for a SIGINT or SIGTERM signal on a 
	//different thread
	go func() {
		//moved side effect out of the log
		sig := <-sigs
		log.Println("received a", sig)
		//make sure that the listener is closed before graceful exit
		l.Close()
	}()
}

func main() {
	//load persistence file 
	err := loadKv()

	if err != nil {
		log.Fatal("Failed to load kv persistence file", err)
	}

	log.Println("starting server on localhost 4000")
	l, err := net.Listen("tcp", ":4000")
	panicIfErr(err)

	handleSysSignals(l)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("error in accept", err)
			break
		}
		log.Println("accepted connection")
		go handle(c)
	}
	//save persistence file
	err = persistKv()

	if err != nil {
		log.Fatal("Failed to persist kv file", err)
	}
}
