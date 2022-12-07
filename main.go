package main

import (
	"flag"
	"log"

	"github.com/peng225/cotton/web"
)

func main() {
	var port int
	var dumpReceivedData bool
	var tls bool
	var serverCrt, serverKey string
	flag.IntVar(&port, "port", 8080, "Port number.")
	flag.BoolVar(&dumpReceivedData, "dump", false, "Dump received data.")
	flag.BoolVar(&tls, "tls", false, "Use TLS.")
	flag.StringVar(&serverCrt, "crt", "", "Server certificate of TLS.")
	flag.StringVar(&serverKey, "key", "", "Server private key of TLS.")
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	if port < 0 || port > 65536 {
		log.Fatalf("Invalid port number: %d", port)
	}

	if tls {
		if serverCrt == "" || serverKey == "" {
			log.Fatal("Both -crt and -key options must be specified to enable TLS.")
		}
	}

	if tls {
		web.StartTLSServer(port, serverCrt, serverKey, dumpReceivedData)
	} else {
		web.StartServer(port, dumpReceivedData)
	}
}
