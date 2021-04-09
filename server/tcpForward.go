package server

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	_ "strings"
	log"github.com/sirupsen/logrus"
)

type tcpForwarderConfig struct {
	localPort string
}

// main serves as the program entry point
func StartForwardServer(destPort uint16) {
	port := "0.0.0.0:8282"
	// create a tcp listener on the given port
	listener, err := net.Listen("tcp4", port)
	if err != nil {
		log.Errorf("failed to create listener, err:", err)
		os.Exit(1)
	}
	log.Infof("forwarder listening on %s\n", listener.Addr())

	// listen for new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("failed to accept connection, err:", err)
			continue
		}
		go handleConnection(conn, destPort)
	}
}

func handleConnection(conn net.Conn, port uint16) {
	newConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		log.Errorf("failed to dial, err:", err)
		return
	}
	reader := bufio.NewReader(conn)
	buf, _ := reader.Peek(1024)
	httpReader := bufio.NewReader(bytes.NewReader(buf))
	_, err = http.ReadRequest(httpReader)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("%v\n", req)
	go flush(conn, reader, newConn)
	go flush(newConn, reader, conn)
}

func flush(src net.Conn, srcReader *bufio.Reader, dst net.Conn) {
	for {
		buf := make([]byte, 1024)
		recvd, err := srcReader.Read(buf)
		//fmt.Printf("%s\n", buf[0:recvd])
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
	}
}
