package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

// main serves as the program entry point
func StartEchoServer() {
	// obtain the port and prefix via program arguments
	port := "127.0.0.1:8181"
	prefix := "test"

	// create a tcp listener on the given port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("failed to create listener, err:", err)
	}
	fmt.Printf("listening on %s, prefix: %s\n", listener.Addr(), prefix)

	// listen for new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("failed to accept connection, err:", err)
			continue
		}

		// pass an accepted connection to a handler goroutine
		go echo(conn, prefix)
	}
}

// handleConnection handles the lifetime of a connection
func echo(conn net.Conn, prefix string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		// read client request data
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				fmt.Println("failed to read data, err:", err)
			}
			return
		}
		fmt.Printf("request: %s", bytes)

		// prepend prefix and send as response
		line := fmt.Sprintf("%s %s", prefix, bytes)
		fmt.Printf("response: %s", line)
		conn.Write([]byte(line))
	}
}
