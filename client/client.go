package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	con, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Printf("error while connecting to server: %s", err)
		return
	}

	defer func(con net.Conn) {
		err := con.Close()
		if err != nil {
			fmt.Printf("error while closing connection: %v\n", err)
		}
	}(con)

	fmt.Println("Connected to the Key value store, please enter your commands")

	clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(con)

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest = strings.TrimSpace(clientRequest)
			if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
				fmt.Printf("failed to send the client request: %v \n", err)
			}
		case io.EOF:
			fmt.Println("client closed the connection")
			return
		default:
			fmt.Printf("client error: %v\n", err)
			return
		}

		// Waiting for the server response
		serverResponse, err := serverReader.ReadString('\n')

		switch err {
		case nil:
			fmt.Println(strings.TrimSpace(serverResponse))
		case io.EOF:
			fmt.Println("server closed the connection")
			return
		default:
			fmt.Printf("server error: %v\n", err)
			return
		}
	}
}
