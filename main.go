package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/peng225/cotton/web"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port number.")
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	if port < 0 || port > 65536 {
		log.Fatalf("Invalid port number: %d", port)
	}

	portStr := strconv.Itoa(port)

	var httpServer http.Server
	http.HandleFunc("/", web.Handler)
	log.Println("Start server.")
	httpServer.Addr = ":" + portStr
	log.Println(httpServer.ListenAndServe())
}
