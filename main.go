package main

import (
	"flag"
	"log"

	"github.com/peng225/cotton/web"
)

func main() {
	var port int
	var dumpReceivedData bool
	flag.IntVar(&port, "port", 8080, "Port number.")
	flag.BoolVar(&dumpReceivedData, "dump", false, "Dump received data.")
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	if port < 0 || port > 65536 {
		log.Fatalf("Invalid port number: %d", port)
	}

	web.StartServer(port, dumpReceivedData)
}
