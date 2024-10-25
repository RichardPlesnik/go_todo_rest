package main

import (
	"flag"
	"log"
)

func main() {
	var port uint

	flag.UintVar(&port, "p", 8080, "port for the server (shorthand)")
	flag.UintVar(&port, "port", 8080, "port for the server")

	flag.Parse()

	log.Println("Starting server")

}
