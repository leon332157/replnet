package client

import (
	"context"
	"time"
	"nhooyr.io/websocket"
)

func StartWS() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
defer cancel()

c, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
if err != nil {
	// ...
}
defer c.Close(websocket.StatusInternalError, "the sky is falling")

c.Close(websocket.StatusNormalClosure, "")
}