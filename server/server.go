package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"keyvaluestore"
	"log"
	"net"
	"strings"
)

var store map[interface{}]interface{}
var result interface{}

func main() {
	store = make(map[interface{}]interface{})

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
		var input keyvaluestore.Command
		rawInput, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			rawInput = strings.TrimSpace(rawInput)
			err := json.Unmarshal([]byte(rawInput), &input)
			if err != nil {
				fmt.Printf("error while unmarshalling: ", err)
			}

			fmt.Println("request received from client: ", input)

		default:
			fmt.Printf("error: %v\n", err)
			return
		}

		finalResult := processInput(input).(string) + "\n"

		// Responding to the client request
		if _, err = con.Write([]byte(finalResult)); err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}

		fmt.Println("Response sent to the client!")

	}
}

func processInput(input keyvaluestore.Command) interface{} {

	switch input.Operation {
	case "SET":
		store[input.Operand1] = input.Operand2
		result = "OK"
	case "GET":
		result = "nil"
		if val, ok := store[input.Operand1]; ok {
			result = val
		}
	case "DELETE":
		result = store[input.Operand1]
		delete(store, input.Operand1)

	default:
		result = "Unsupported command or Wrong Syntax! Please try again"
	}

	return result
}
