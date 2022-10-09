package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	var address = flag.String("address", "", "address to listen to")

	flag.Parse()
	if *address == "" {
		flag.PrintDefaults()
		log.Fatalln("No address specified.")
	}

	fmt.Println("Listening on address", *address)
	startNetDialWrite(*address)
}

func startNetDialWrite(addr string) {
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

	buf := bytes.NewBuffer([]byte("add 10"))
	buf.WriteRune(';')

	if _, err := io.Copy(conn, buf); err != nil {
		log.Fatalln(err)
	}

	time.Sleep(1 * time.Second)

	buf.Reset()
	buf.Write([]byte("close"))
	buf.WriteRune(';')

	if _, err := io.Copy(conn, buf); err != nil {
		log.Fatalln(err)
	}
}
