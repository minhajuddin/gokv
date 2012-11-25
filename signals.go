package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

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
