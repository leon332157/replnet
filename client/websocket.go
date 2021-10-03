package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
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
	MaxIdleConns:          5,
	IdleConnTimeout:       30 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
var httpClient = http.Client{Transport: transport}

func keepAlive(c *websocket.Conn) {
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

func connectWS(remoteUrl string, remotePort uint16, timeout time.Duration) {
	/*  if remoteUrl == "" {
		log.Fatalf("[Websocket Client] remoteUrl is empty")
		return
	}*/
	remoteUrl = strings.TrimRight(remoteUrl, "/")
	log.Debugf("[Websocket Client] Connecting to %v", remoteUrl)
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	//defer cancel()
	c, _, err := websocket.Dial(ctx, fmt.Sprintf("%s/__ws?remoteAppPort=%v", remoteUrl, remotePort), &websocket.DialOptions{HTTPClient: &httpClient})
	if err != nil {
		log.Fatalf("[Websocket Client] Dial failed: %s", err)
	}
	//c.Write(ctx,websocket.MessageText,[]byte("Test"))
	//defer c.Close(websocket.StatusInternalError, "the sky is falling")
	go func() {
		for {
			msgtype, data, err := c.Read(context.Background())
			log.Debugf("[WS Client] type: %s data: %s err: %v", msgtype, data, err)
			if err != nil {
				c.Close(websocket.StatusInternalError, err.Error())
			}
			//time.Sleep(1000*time.Millisecond)
		}
	}()
	//go keepAlive(c)
	//c.Close(websocket.StatusNormalClosure, "")
}
