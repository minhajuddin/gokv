//gokv is a basic key value store which is persisted to disk in a json serialized form
//it supports 
//	+ storing a key value pair
//	+ retreiving the value for a given key
//	+ deleting a key value pair
//	+ retreiving a list of keys with a given prefix

package main

import (
	"log"
	"net"
	"os"
	"path"
	"sync"
)

var (
	//file where the serialized json is persisted
	kvFile = path.Join(os.Getenv("HOME"), ".gokv.json")

	//mutex used to lock/unlock access to the kv store
	mutex = &sync.Mutex{}
	//in-memory key value store
	kv map[string]interface{}
)

func main() {
	//load persistence file 
	err := loadKv()

	if err != nil {
		log.Fatal("Failed to load kv persistence file", err)
	}

	log.Println("starting server on localhost 4000")
	l, err := net.Listen("tcp", ":4000")

	if err != nil {
		log.Fatalln("failed to start the server", err)
	}

	handleSysSignals(l)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("error in accept", err)
			break
		}
		log.Println("accepted connection", c.RemoteAddr())
		go handle(c)
	}
	//save persistence file
	err = persistKv()

	if err != nil {
		log.Fatal("Failed to persist kv file", err)
	}
}
