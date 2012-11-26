package main

import (
	"bufio"
	"log"
	"net"
)

//core handler
func handle(c net.Conn) {
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
		//execute the command
		cmd.Exec(w)
	}
}
