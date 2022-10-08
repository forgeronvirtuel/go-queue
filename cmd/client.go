package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	var address = flag.String("address", "", "address to listen to")

	flag.Parse()
	if *address == "" {
		flag.PrintDefaults()
		log.Fatalln("No address specified.")
	}

	fmt.Println("Listening on address", *address)
	startNetDial(*address)
}

func startNetDial(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		log.Println("Closing connection")
		if err := conn.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	buf := make([]byte, 1024)
	if _, err := conn.Read(buf); err != nil {
		log.Fatalln(err)
	}

	log.Println(string(buf))
}
