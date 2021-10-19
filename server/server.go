package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type operations interface {
	get() interface{}
	delete() (interface{}, error)
	set(interface{}) error
}

type Command struct {
	operation string
	operand1  interface{}
	operand2  interface{}
}

type ConnectionContext struct {
	connection net.Conn
	owner      string
}

var activeConns = make(chan ConnectionContext)
var connMap = make(map[net.Conn]string)
var deadConns = make(chan net.Conn, 10)
var results = make(chan interface{})
var store map[interface{}]interface{}
var result interface{}

func main() {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	store = make(map[interface{}]interface{})

	service := fmt.Sprintf("0.0.0.0:%s", port)
	listener, err := net.Listen("tcp", service)
	if err != nil {
		fmt.Errorf("server: listen: %s", err)
	}

	fmt.Printf("server: listening on port %s", port)

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(con)
	}
}

func handleConnection(con net.Conn) {
	defer con.Close()

	clientReader := bufio.NewReader(con)

	for {
		// Waiting for the client request
		incoming, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest := strings.TrimSpace(incoming)
			if clientRequest == ":QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			} else {
				log.Println("Here's the request", clientRequest)
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}

		//var cmd Command
		//
		//if incoming != "" {
		//	cmd.operation = strings.Split(incoming, " ")[0]
		//	cmd.operand1 = strings.Split(incoming, " ")[1]
		//	cmd.operand2 = strings.Split(incoming, " ")[2]
		//	switch cmd.operation {
		//	case "SET":
		//		store[cmd.operand1] = cmd.operand2
		//		result = "OK"
		//	case "GET":
		//		result = store[cmd.operand1]
		//	case "DELETE":
		//		result = store[cmd.operand1]
		//		delete(store, cmd.operand1)
		//	}

		// Responding to the client request
		if _, err = con.Write([]byte(incoming)); err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}
	}

}
