package main

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
)

//TODO: write a few tests to make sure that cmd parsing is fine
var cmdrx = regexp.MustCompile("(GET|SET|LIST|DELETE) ?([a-zA-Z0-9_-]+)? ?(.*)")

type command struct {
	ctype string
	key   string
	value interface{}
}

//refer to the protocol in the README
func (self *command) Exec(w *bufio.Writer) {
	defer w.Flush()
	switch self.ctype {
	case "GET":
		val, ok := getValue(self.key)
		if ok {
			w.WriteString(fmt.Sprintln(val))
		} else {
			w.WriteString("<NULL>\n")
		}
	case "SET":
		setValue(self.key, self.value)
	case "LIST":
		for _, k := range listKeys(self.key) {
			w.WriteString(k + "\n")
		}
	case "DELETE":
		deleteKey(self.key)
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
