package server

import (
	"bufio"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	_ "net/http"
	"os"
	"github.com/valyala/fasthttp"
)

type tcpForwarderConfig struct {
	localPort string
}

func UNUSED(x ...interface{}) {}

// main serves as the program entry point
func StartForwardServer(destPort uint16) {
	// create a tcp listener on port assigned by kernel
	listener, err := net.Listen("tcp4", "0.0.0.0:0")
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
	go flush(conn, newConn)
	go flush(newConn, conn)
}

func flush(src net.Conn, dst net.Conn) {
	for {
		buf := make([]byte, 1024000)
		recvd, err := src.Read(buf)
		//fmt.Printf("%s\n", buf[0:recvd])
		if err != nil {
			log.Errorf("error reading %v %v\n", src.RemoteAddr(), err)
			dst.Close()
			src.Close()
			return
		}
		fmt.Printf("%s\n",buf[0:100])
		if bytes.Contains(buf[0:100], []byte("HOST")) { // if this is a HTTP request

			httpReader := bufio.NewReader(bytes.NewReader(buf[0:8192]))
			newReq := fasthttp.AcquireRequest()
			err = newReq.Read(httpReader)
			if err != nil {
				log.Error(err)
			}
			//UNUSED(req)
			//log.Debugf("request: %s\n", newReq)
		}
		sent, err := dst.Write(buf[0:recvd])
		log.Debugf("flushed %v bytes to %v\n", sent, dst.RemoteAddr())
		if err != nil {
			log.Errorf("error sending to %v %v\n", dst.RemoteAddr(), err)
			dst.Close()
			src.Close()
			return
		}
		buf = nil
	}
}
