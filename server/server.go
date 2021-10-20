package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type command struct {
	operation string
	operand1  string
	operand2  interface{}
}

var store map[string]interface{}
var result interface{}

func main() {
	store = make(map[string]interface{})

	service := fmt.Sprintf("0.0.0.0:8080")
	listener, err := net.Listen("tcp", service)
	if err != nil {
		fmt.Printf("error while starting server: %s", err)
		return
	}

	fmt.Printf("server: listening on port 8080")

	for {
		con, err := listener.Accept()
		if err != nil {
			fmt.Printf("error while connecting to client: %v\n", err)
			continue
		}

		go handleConnection(con)
	}
}

func handleConnection(con net.Conn) {
	defer func(con net.Conn) {
		err := con.Close()
		if err != nil {
			fmt.Printf("error while closing connection: %v\n", err)
		}
	}(con)

	clientReader := bufio.NewReader(con)

	for {
		// Waiting for the client
		input, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			input = strings.TrimSpace(input)
			if input == "quit" {
				fmt.Println("closing the connection as requested by client")
				return
			} else {
				fmt.Println("request received from client: ", input)
			}
		case io.EOF:
			fmt.Println("client closed the connection by terminating the process")
			return
		default:
			fmt.Printf("error: %v\n", err)
			return
		}

		if input != "" {
			finalResult := processInput(input).(string) + "\n"

			// Responding to the client request
			if _, err = con.Write([]byte(finalResult)); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}

			fmt.Println("Response sent to the client!")
		}

	}
}

func processInput(input string) interface{} {
	var cmd command
	cmd.operation = strings.Split(input, " ")[0]
	cmd.operand1 = strings.Split(input, " ")[1]

	switch cmd.operation {
	case "SET":
		cmd.operand2 = strings.Split(input, " ")[2]
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
		result = "Unsupported command or Wrong Syntax! Please try again"
	}

	return result
}
