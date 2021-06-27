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

func heartbeat(ctx context.Context, c *websocket.Conn, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Debugln("done")
			return
		case <-t.C:
		}
		err := c.Ping(ctx)
		if err != nil {
			log.Debugln(err)
		} else {
			log.Debugln("Ping!")
		}

		t.Reset(time.Second)
	}
}
func StartWS() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	c, _, err := websocket.Dial(ctx, "ws://localhost:7070/__ws", &websocket.DialOptions{HTTPClient: &httpClient})
	if err != nil {
		log.Fatalf("[Websocket Client] Dial failed: %s", err)
	}
	c.Write(ctx,websocket.MessageText,[]byte("Test"))
	//defer c.Close(websocket.StatusInternalError, "the sky is falling")
	/*err = c.Ping(ctx)
	if err != nil {
		log.Debugln(err)
	} else {
		log.Debugln("PING")
	}*/
	go func() {
		for {
			timeout, _ := context.WithTimeout(context.Background(), 5*time.Second)
			err := c.Ping(timeout)
			//c.Read(ctx)
			//c.Write(timeout,websocket.MessageText,[]byte("PING"))
			log.Debugln("[Websocket Client] Keep alive")
			if err != nil {
				log.Debugf("[Websocket Client] Keep alive err: %s\n", err)

			}
			time.Sleep(1 * time.Second)
		}
	}()
	//hb := context.TODO()
	//go heartbeat(ctx, c, time.Second)
	//c.Close(websocket.StatusNormalClosure, "")
}
