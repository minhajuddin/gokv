package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
)

//TODO: write a few tests to make sure that cmd parsing is fine
var cmdrx = regexp.MustCompile("(GET|STORE|LIST|DELETE) ([a-zA-Z0-9_-]+) ?(.*)")

//core handler
//if the input is just one line it tries to get the
//value for the given key
//if it is more than one line, it assumes that the 
//first line is a key and the rest is the value
//it tries to store the data in this scenario
func handle(c net.Conn) {
	//TODO: should be an infiinte loop?
	remoteAddr := c.RemoteAddr()
	defer log.Println("connection closed for ", remoteAddr)
	defer c.Close()
	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)
	for {
		rawline, err := r.ReadString('\n')

		//connection has probably closed
		if err != nil {
			log.Println(err)
			return
		}

		cmd, err := parseCommand(rawline)

		if err != nil {
			w.WriteString("<INVALID COMMAND>\n")
			w.Flush()
			//n, err := w.WriteString("<INVALID COMMAND>")
			log.Println("Invalid command:", rawline, err, "from", remoteAddr)
			continue
		}

		log.Println("processing", cmd, "from", remoteAddr)
		cmd.Exec(*w)
		////parse
		//buf := make([]byte, bufsize)
		//nr, _ := io.ReadFull(c, buf)
		//key := strings.Trim(string(buf[:nr]), "\n")
		//idx := strings.Index(key, "\n")
		////if it has a newline, the first line is the key
		//if idx > -1 {
		//setValue(key[:idx], strings.Trim(key[idx:], "\n"))
		//return
		//}
		//if v, ok := getValue(key); ok {
		//fmt.Fprintln(c, v)
		//return
		//}
		//log.Printf("key for '%s' not found for '%v'\n", key, remoteAddr)
		//fmt.Fprintln(c, "<NULL>")
	}
}

type command struct {
	ctype string
	key   string
	value interface{}
}

func (self *command) Exec(w bufio.Writer) {
	defer w.Flush()
	switch self.ctype {
	case "GET":
		val, ok := getValue(self.key)
		if ok {
			w.WriteString(fmt.Sprintln(val))
		} else {
			w.WriteString("<NULL>\n")
		}
	}
}

func parseCommand(rawline string) (*command, error) {
	if !cmdrx.MatchString(rawline) {
		return nil, errors.New("Invalid command")
	}
	tokens := cmdrx.FindStringSubmatch(rawline)

	cmd := &command{
		ctype: tokens[1],
		key:   tokens[2],
		value: tokens[3],
	}

	return cmd, nil
}
