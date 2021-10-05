package server

import (
	"context"
	//"time"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	websocket "nhooyr.io/websocket"
)

func UNUSED(x ...interface{}) {
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	log.Debugln("[WS handler] recvd")
	fmt.Println(r.URL.Query())
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Errorf("[Websocket Handler] Accept from %v err: %v\n", r.RemoteAddr, err)
		return
	} else {
		log.Debugf("[Websocket Handler] Accepted from %v", r.RemoteAddr)
	}
	go func() {
		for {
			msgtype, data, err := c.Read(context.Background())
			log.Debugf("[WS handler] type: %s data: %s err: %v", msgtype, data, err)
			if err != nil {
				c.Close(websocket.StatusInternalError, fmt.Sprintf("[Websocker Handler] Read err: %v", err))
				return			
			}
			//time.Sleep(1000*time.Millisecond)
		}
	}()
	//c.Close(websocket.StatusNormalClosure, "")
}
