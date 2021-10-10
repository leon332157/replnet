package server

import (
	"context"
	//"time"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	websocket "nhooyr.io/websocket"
	"strconv"
	//"strings"
)

func handleWS(w http.ResponseWriter, r *http.Request) {
	log.Debugln("[WS handler] recvd")
	stringPort := r.URL.Query().Get("remoteAppPort")
	if len(stringPort) == 0 {
		return
	}
	intPort, err := strconv.ParseUint(stringPort, 10, 16)
	if err != nil {
		return
	}
	port := uint16(intPort)
	log.Debugf("[WS handler] remoteAppPort: %d", port)
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Errorf("[Websocket Handler] Accept from %v err: %v\n", r.RemoteAddr, err)
		return
	} else {
		log.Debugf("[Websocket Handler] Accepted from %v", r.RemoteAddr)
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		log.Errorf("[Websocket Handler] Dial to %v err: %v\n", port, err)
		c.Close(websocket.StatusInternalError, err.Error())
		return
	}

	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				log.Errorf("[Websocket Handler] Read socket from %v err: %v\n", r.RemoteAddr, err)
				conn.Close()
				c.Close(websocket.StatusInternalError, err.Error())
				return
			}
			c.Write(context.Background(), websocket.MessageBinary, buf[:n])
		}
	}()

	go func() {
		for {
			msgtype, data, err := c.Read(context.Background())
			log.Debugf("[WS handler] type: %s data: %s err: %v", msgtype, data, err)
			if err != nil {
				c.Close(websocket.StatusInternalError, fmt.Sprintf("[Websocker Handler] Read err: %v", err))
				return
			}
			if msgtype == websocket.MessageBinary {
				conn.Write(data)
			}
			//time.Sleep(1000*time.Millisecond)
		}
	}()
	//c.Close(websocket.StatusNormalClosure, "")
}
