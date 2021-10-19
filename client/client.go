package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type ClientMetadata struct {
	Owner          string
	ChatServerHost string
	ChatServerPort string
	OwnerEmail     string
}

type Command struct {
	operation string
	operand1  interface{}
	operand2  interface{}
}

var clientMD ClientMetadata
var conn net.Conn
var connectionC = make(chan net.Conn)
var commands = make(chan Command)

func main() {
	con, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(con)

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
				log.Printf("failed to send the client request: %v \n", err)
			}
		case io.EOF:
			log.Println("client closed the connection")
			return
		default:
			log.Printf("client error: %v\n", err)
			return
		}

		// Waiting for the server response
		serverResponse, err := serverReader.ReadString('\n')

		switch err {
		case nil:
			log.Println(strings.TrimSpace(serverResponse))
		case io.EOF:
			log.Println("server closed the connection")
			return
		default:
			log.Printf("server error: %v\n", err)
			return
		}
	}
}
