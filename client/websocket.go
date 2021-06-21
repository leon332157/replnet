package client

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func StartWS() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, _, err := websocket.Dial(ctx, "ws://localhost:7070/__ws", nil)
	if err != nil {
		log.Errorf("[Websocket Client] Dial failed: %s", err)
	}

	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	err = wsjson.Write(ctx, c, "hi")
	c.Ping(ctx)
	if err != nil {
		fmt.Println(err)
	}

	c.Close(websocket.StatusNormalClosure, "")
}
