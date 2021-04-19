package server

import (
	"bufio"
	"bytes"
	"fmt"
	_ "io"
	"net"
	_ "net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func UNUSED(x ...interface{}) {}

const REMOTE_PORT uint16 = 8383

// main serves as the program entry point
func StartForwardServer(destPort uint16) {
	// create a tcp listener on port assigned by kernel
	var addr *net.TCPAddr
	addr, _ = net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%v", REMOTE_PORT))
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Errorf("failed to create listener, err:%v", err)
		os.Exit(1)
	}
	log.Infof("forwarder listening on %s\n", listener.Addr())
	addr, _ = net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%v", destPort))
	localConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Errorf("failed to dial, err:%v", err)
		return
	}
	localConn.SetKeepAlive(true)
	localConn.SetKeepAlivePeriod(5 * time.Second)
	// listen for new connections
	for {
		var conn *net.TCPConn
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Errorf("failed to accept connection, err:%v", err)
			continue
		}
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(5 * time.Second)
		//go io.Copy(conn, localConn)
		//go io.Copy(localConn, conn)
		go flush(conn, localConn)
		go flush(localConn, conn) // Use io.Copy eventually
	}
}

// TODO: have two different flush functions for each direction
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
		//fmt.Printf("%s\n", buf[0:100])
		if bytes.Contains(buf[0:100], []byte("HOST")) { // if this is a HTTP request
			httpReader := bufio.NewReader(bytes.NewReader(buf[0:2048])) // read 2048
			newReq := fasthttp.AcquireRequest()
			err = newReq.Read(httpReader)
			if err != nil {
				log.Error(err)
			}
			//UNUSED(req)
			//log.Debugf("request: %s\n", newReq)
		}
		if len(buf[0:recvd]) != 0 {
			sent, err := dst.Write(buf[0:recvd])
			log.Debugf("flushed %v bytes to %v\n", sent, dst.RemoteAddr())
			if err != nil {
				log.Errorf("error sending to %v %v\n", dst.RemoteAddr(), err)
				dst.Close()
				src.Close()
				return
			}
		}
		buf = nil
	}
}
