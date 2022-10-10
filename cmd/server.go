package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/forgeronvirtuel/mqueue/command"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
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
	var connCounter int

	if err != nil {
		log.Fatalln(err)
	}

	// Create the queue
	//datach := make(chan []byte)
	//queueHandler := NewQueueHandler(datach)

	// Securely close server
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
		go handleConn(conn, connCounter)
		connCounter++
	}
}

func handleConn(conn net.Conn, counter int) {
	cmds := &command.Manager{}
	logger := log.New(os.Stdout, fmt.Sprintf("[%d] ", counter), log.Ldate|log.Ltime|log.Lshortfile)
	cmds.Add(command.NewItem("add"), func(params []string) error {
		if len(params) != 1 {
			return fmt.Errorf("Detected %d params, wants %d params", len(params), 1)
		}
		length, err := strconv.Atoi(params[0])
		if err != nil {
			return fmt.Errorf("cannot convert <length> param to int: %v", err)
		}
		logger.Printf("add %d detected", length)
		return nil
	})
	receiveData(conn, cmds, logger)
}

const buf_size = 1024

const CLOSECMD = "close"

// todo protect against multicommands
func receiveData(conn net.Conn, cmds *command.Manager, logger *log.Logger) {
	// Receive part
	buf := bytes.NewBuffer(make([]byte, 0, buf_size))
	buftmp := make([]byte, buf_size)

	// adding end mecanism
	var endConn bool
	cmds.Add(command.NewItem(CLOSECMD), func(_ []string) error {
		endConn = true
		logger.Println("Closing client connection")
		return nil
	})

	// reading data
	for !endConn {
		var done bool
		for !done {
			nread, err := conn.Read(buftmp)
			if err == io.EOF {
				logger.Println("Client closed connection")
				done = true
				endConn = true
				continue
			}
			if err != nil {
				logger.Println("Error at reading", err)
			}
			ncpy, err := buf.Write(buftmp[0:nread])
			if ncpy != nread {
				logger.Println("Cannot copy buftmp to buf")
			}
			if err != nil {
				logger.Println(err)
			}

			// Securized the buffer size
			if nread+buf.Len() >= buf_size {
				logger.Println("Total size of buffer is outcapacited")
			}

			if buftmp[nread-1] == ';' {
				done = true
			}
		}

		// In case the connection has stopped
		if endConn {
			continue
		}

		var err error
		length := buf.Len()
		if _, err = buf.Read(buftmp); err != nil {
			logger.Println(err)
		}
		buf.Reset()

		action := cmds.Parse(buftmp[:length-1])
		if action == nil {
			logger.Printf("No action for command %s", string(buftmp))
		}
		if err := action(); err != nil {
			logger.Printf("Error while treating command: %v", err)
		}
	}
	logger.Println("Client connection closed")
}

//func NewQueueHandler(data chan []byte) *QueueHandler {
//	return &QueueHandler{data: data}
//}

//type QueueHandler struct {
//	data chan []byte
//	done bool
//	mq   queue.Queue
//}
//
//func (q *QueueHandler) run() {
//	for !q.done {
//		q.mq.Produce(<-q.data)
//	}
//}
//
//func (q *QueueHandler) close() {
//	q.done = true
//}
