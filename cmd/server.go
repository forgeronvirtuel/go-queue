package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	var address = flag.String("address", "", "address to listen to")

	flag.Parse()
	if *address == "" {
		flag.PrintDefaults()
		log.Fatalln("No address specified.")
	}

	fmt.Println("Listening on address", *address)
	startNetServer(*address)
}

func startNetServer(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Closing server")
		if err := l.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	for {
		// Listen for an incoming connection
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	var n int
	var err error
	msg := []byte("Hello, world !")
	for n != len(msg) {
		if n, err = conn.Write(msg); err != nil {
			log.Fatalln(err)
		}
	}
	if err = conn.Close(); err != nil {
		log.Fatalln(err)
	}
}
