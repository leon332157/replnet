package client

import (
	"context"
	log "github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
	"time"
)

func heartbeat(ctx context.Context, c *websocket.Conn, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	for {
		log.Debugln("a")
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
	c, _, err := websocket.Dial(ctx, "ws://localhost:7070/__ws", nil)
	if err != nil {
		log.Fatalf("[Websocket Client] Dial failed: %s", err)
	}

	//defer c.Close(websocket.StatusInternalError, "the sky is falling")
	/*go func() {
		for {
			timeout, _ := context.WithTimeout(context.Background(), 5*time.Second)
			err := c.Ping(timeout)
			log.Debugln("[Websocket Client] Keep alive")
			if err != nil {
				log.Debugf("[Websocket Client] Keep alive err: %s\n", err)
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()*/
	hb := context.TODO()
	go heartbeat(hb, c, time.Second)
	//c.Close(websocket.StatusNormalClosure, "")
}
