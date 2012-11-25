package main

import (
	"io"
	"log"
	"net"
	"strings"
	"fmt"
)

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
