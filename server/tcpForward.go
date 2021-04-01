package server

import (
	"fmt"
	"net"
	"os"
)

// main serves as the program entry point
func StartForwardServer() {
	port := "0.0.0.0:8282"
	// create a tcp listener on the given port
	listener, err := net.Listen("tcp4", port)
	if err != nil {
		fmt.Println("failed to create listener, err:", err)
		os.Exit(1)
	}
	fmt.Printf("listening on %s\n", listener.Addr())

	// listen for new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("failed to accept connection, err:", err)
			continue
		}
		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	newConn, err := net.Dial("tcp", "127.0.0.1:8181")
	if err != nil {
		fmt.Println("failed to dial, err:", err)
		return
	}
	//defer newConn.Close()
	//defer conn.Close()
	go flush(conn, newConn)
	go flush(newConn, conn)
}

func flush(src net.Conn, dst net.Conn) {
	for {
		buf := make([]byte, 1048576)
		recvd, err := src.Read(buf)
		if err != nil {
			fmt.Printf("error %v %v\n", src.RemoteAddr(), err)
			dst.Close()
			src.Close()
			return
		}
		sent, err := dst.Write(buf[0:recvd])
		fmt.Printf("flushed %v bytes to %v\n", sent, dst.RemoteAddr())
		if err != nil {
			fmt.Printf("error sending to %v %v\n", dst.RemoteAddr(), err)
			dst.Close()
			src.Close()
			return
		}
		buf = nil
		//runtime.GC()
	}
}
