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
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Errorf("[Websocker Handler] Accept err: %v\n", err)
	}
	go func() {
		for {
			msgtype, data, err := c.Read(context.Background())
			log.Debugf("[WS handler] type: %s data: %s err: %v", msgtype, data, err)
			if err != nil {
				c.Close(websocket.StatusInternalError, fmt.Sprintf("[Websocker Handler] Read err: %v", err))
			}
			//time.Sleep(1000*time.Millisecond)
		}
	}()
	//c.Close(websocket.StatusNormalClosure, "")
}
