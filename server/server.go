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

type Command struct {
	operation string
	operand1  string
	operand2  interface{}
}

var store map[string]interface{}
var result interface{}

func main() {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	store = make(map[string]interface{})

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
		// Waiting for the client
		incoming, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest := strings.TrimSpace(incoming)
			if clientRequest == ":QUIT" {
				log.Println("Closing the connection as requested by client")
				return
			} else {
				log.Println("Request received from client: ", clientRequest)
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}

		var cmd Command

		if incoming != "" {
			cmd.operation = strings.Split(incoming, " ")[0]
			cmd.operand1 = strings.Split(incoming, " ")[1]

			fmt.Println(cmd)
			switch cmd.operation {
			case "SET":
				cmd.operand2 = strings.Split(incoming, " ")[2]
				store[cmd.operand1] = strings.TrimSpace(cmd.operand2.(string))
				result = "OK"
			case "GET":
				result = "nil"
				if val, ok := store[strings.TrimSpace(cmd.operand1)]; ok {
					result = val
				}
			case "DELETE":
				result = store[strings.TrimSpace(cmd.operand1)]
				delete(store, strings.TrimSpace(cmd.operand1))

			default:
				result = "Unsupported command!"
			}

			finalResult := result.(string) + "\n"

			// Responding to the client request
			if _, err = con.Write([]byte(finalResult)); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}

			fmt.Println("Response sent to the client!")
		}

	}
}
