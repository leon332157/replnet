package client

import (
	"context"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
)

var transport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
var httpClient = http.Client{Transport: transport}

func keepAlive(c websocket.Conn) {
	for {
		timeout, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := c.Ping(timeout)
		//c.Write(timeout,websocket.MessageText,[]byte("PING"))
		log.Debugln("[Websocket Client] Keep alive")
		if err != nil {
			log.Debugf("[Websocket Client] Keep alive err: %s\n", err)

		}
		time.Sleep(1 * time.Second)
	}
}
func StartWS() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	c, _, err := websocket.Dial(ctx, "ws://localhost:7777/__ws", &websocket.DialOptions{HTTPClient: &httpClient})
	if err != nil {
		log.Fatalf("[Websocket Client] Dial failed: %s", err)
	}
	//c.Write(ctx,websocket.MessageText,[]byte("Test"))
	//defer c.Close(websocket.StatusInternalError, "the sky is falling")
	go func() {
		for {
			_, data, err := c.Read(context.Background())
			log.Debugf("[WS handler] data: %s, err: %v", data, err)
			if err != nil {
				break
			}
			//time.Sleep(1000*time.Millisecond)
		}
	}()
	//c.Close(websocket.StatusNormalClosure, "")
}
