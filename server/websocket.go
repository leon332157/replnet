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
	"io"
	//"strings"
)

func sockToWs(ws *websocket.Conn, sock net.Conn) {
for {
			buf := make([]byte, 256)
			n, err := sock.Read(buf)
			if err != nil {
				if err != io.EOF{
				log.Errorf("[Websocket Handler] Read socket from %v err: %v\n", sock.RemoteAddr, err)
				//sock.Close()
				//ws.Close(websocket.StatusInternalError, err.Error())
				ws.Write(context.Background(), websocket.MessageBinary, []byte("baf"))
				}
				return
			}
			ws.Write(context.Background(), websocket.MessageBinary, buf[:n])
		}
}

func wsToSock(ws *websocket.Conn, sock net.Conn) {
	for {
			_, data, err := ws.Read(context.Background())
			log.Debugf("[WS handler] data: %s err: %v", data, err)
			if err != nil {
				//c.Close(websocket.StatusInternalError, fmt.Sprintf("[Websocker Handler] Read err: %v", err))
				return
			}
			sock.Write(data)
			//time.Sleep(1000*time.Millisecond)
		}
}

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
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Errorf("[Websocket Handler] Accept from %v err: %v\n", r.RemoteAddr, err)
		return
	} else {
		log.Debugf("[Websocket Handler] Accepted from %v", r.RemoteAddr)
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		log.Errorf("[Websocket Handler] Dial to %v err: %v\n", port, err)
		ws.Close(websocket.StatusInternalError, err.Error())
		return
	}
	//io.Copy(ws,conn)
	//io.Copy(conn,ws)
	go sockToWs(ws, conn)
    go wsToSock(ws, conn)
	//c.Close(websocket.StatusNormalClosure, "")
}
