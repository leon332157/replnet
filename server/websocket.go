package server

import (
	"context"
	"time"
	log "github.com/sirupsen/logrus"
	"net/http"
	"nhooyr.io/websocket"
)

func UNUSED(x ...interface{}) {
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Errorf("[Websocker Handler] Accept err: %v\n", err)
	}
	
	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	/*go func() {
		for {
			timeout, _ := context.WithTimeout(context.Background(), 5*time.Second)
			err := c.Ping(timeout)
			log.Debugln("[Websocket Handler] Keep alive")
			if err != nil {
				log.Debugf("[Websocket Handler] Keep alive err: %s\n", err)
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()*/
	go func() {
		for {
			_, data, err := c.Read(context.Background())
			log.Debugf("[WS handler] data: %s, err: %v", data, err)
			if err != nil {
				//break
			}
			time.Sleep(1000*time.Millisecond)
		}
	}()
	//c.Close(websocket.StatusNormalClosure, "")
}
