package server

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	websocket "nhooyr.io/websocket"
	"strconv"
)

// flush socket data to websocket
func sockToWs(ws *websocket.Conn, sock net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := sock.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Debugf("[Websocket Handler] sock read EOF")
				ws.Write(context.Background(), websocket.MessageBinary, buf[:n])
				ws.Close(websocket.StatusBadGateway, "EOF")
				sock.Close()
			} else {
				log.Errorf("[Websocket Handler] read sock from %v err: %v\n", sock.RemoteAddr(), err)
				sock.Close()
				ws.Close(websocket.StatusInternalError, err.Error())
			}
			return
		}
		ws.Write(context.Background(), websocket.MessageBinary, buf[:n])
	}
}

// flush websocket data to socket
func wsToSock(ws *websocket.Conn, sock net.Conn) {
	for {
		_, data, err := ws.Read(context.Background())
		log.Debugf("[Websocket handler] data: %s err: %v", data, err)
		if err != nil {
			//c.Close(websocket.StatusInternalError, fmt.Sprintf("[Websocker Handler] Read err: %v", err))
			return
		}
		n, err := sock.Write(data)
		if err != nil {
			ws.Close(websocket.StatusInternalError, err.Error())
			sock.Close()
			return
		}
		log.Debugf("[Websocket handler] flushed %d bytes to sock", n)
	}
}

func handleControl(ws *websocket.Conn) {
	for {
		msgtype, data, err := ws.Read(context.Background())
		if err != nil {
			log.Errorf("[Websocket handler] read err: %v", err)
			return
		}
		log.Debugf("[Websocket handler]type: %s data: %s", msgtype, data)
	}
}

// Websocket handler
func handleWS(w http.ResponseWriter, r *http.Request) {
	controlFlag := false
	log.Debugln("[Websocket handler] recvd")
	stringPort := r.URL.Query().Get("remoteAppPort")
	if len(stringPort) == 0 {
		controlFlag = true
	}
	intPort, err := strconv.ParseUint(stringPort, 10, 16)
	if err != nil {
		return
	}
	port := uint16(intPort)
	log.Debugf("[WS handler] remoteAppPort: %d", port)
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Errorf("[Websocket Handler] Accept from %v err: %v\n", r.RemoteAddr, err)
		return
	} else {
		log.Debugf("[Websocket Handler] Accepted from %v", r.RemoteAddr)
	}
	if controlFlag {
		go handleControl(ws)
		//ws.Close(websocket.StatusNormalClosure, "Control connection")
		return
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		log.Errorf("[Websocket Handler] Dial to %v err: %v\n", port, err)
		ws.Close(websocket.StatusInternalError, err.Error())
		return
	}
	go sockToWs(ws, conn)
	go wsToSock(ws, conn)
	//c.Close(websocket.StatusNormalClosure, "")
}
