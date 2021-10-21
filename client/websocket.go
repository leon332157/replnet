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

var (
	transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 3 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          5,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	httpClient = http.Client{Transport: transport}
)

func wsToSock(ws *websocket.Conn, sock net.Conn) {
	for {
		defer sock.Close()
		_, data, err := ws.Read(context.Background())
		log.Debugf("[Websocket Client] data: %s err: %v", data, err)
		if err != nil {
			log.Error("[Websocket Client] read from remote error: ", err)
			ws.Close(websocket.StatusInternalError, err.Error())
			return
		}
		written, err := sock.Write(data)
		if err != nil {
			log.Debugf("[Websocket Client] Write failed: %v", err)
			return
		} else {
			log.Debugf("[Websocket Client] flushed %v to sock", written)
		}

	}
}

func sockToWs(ws *websocket.Conn, sock net.Conn) {
	defer sock.Close()
	for {
		buf := make([]byte, 1024)
		recvd, err := sock.Read(buf)
		if err != nil {
			log.Debugf("[Websocket Client] Read failed: %v", err)
			return
		}
		err = ws.Write(context.Background(), websocket.MessageBinary, buf[:recvd])
		if err != nil {
			log.Debugf("[Websocket Client] Write failed: %v", err)
			return
		} else {
			log.Debugf("[Websocket Client] flushed %v to channel", recvd)
		}
	}
}

// main entry for websocket
func startWS(remoteUrl string, remotePort uint16, localPort uint16, timeout time.Duration) {
	/*  if remoteUrl == "" {
		log.Fatalf("[Websocket Client] remoteUrl is empty")
		return
	}*/
	remoteUrl = strings.TrimRight(remoteUrl, "/")
	log.Infof("[Websocket Client] Connecting to %v", remoteUrl)
	ctx := context.Background() //context.WithTimeout(context.Background(), timeout)
	_, _, err := websocket.Dial(ctx, fmt.Sprintf("%s/__ws?remoteAppPort=%v", remoteUrl, remotePort), &websocket.DialOptions{HTTPClient: &httpClient})
	//TODO: add timeout
	//TODO: control channel
	if err != nil {
		log.Fatalf("[Websocket Client] Dial failed: %s", err)
	}
	log.Infof("[Websocket Client] Connected to %v", remoteUrl)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", localPort))
	if err != nil {
		log.Debugf("[Websocket Client] Listen failed: %v", err)
	}
	log.Debugf("[Websocket Client] Local listener created on %v", listener.Addr())

	for {
		sock, err := listener.Accept()
		if err != nil {
			log.Debugf("[Websocket Client] Accept failed: %v", err)
		} else {
			log.Debugf("[Websocket Client] Accepted from: %v", sock.RemoteAddr())
		}
		ws, _, err := websocket.Dial(ctx, fmt.Sprintf("%s/__ws?remoteAppPort=%v", remoteUrl, remotePort), &websocket.DialOptions{HTTPClient: &httpClient})
		if err != nil {
			log.Fatalf("[Websocket Client] Dial failed: %s", err)
		}
		go wsToSock(ws, sock)
		go sockToWs(ws, sock)

	}
	//c.Write(ctx,websocket.MessageText,[]byte("Test"))
	//defer c.Close(websocket.StatusInternalError, "the sky is falling")

	/*go func() {
		for {
			msgtype, data, err := c.Read(context.Background())
			log.Debugf("[WS Client] type: %s data: %s err: %v", msgtype, data, err)
			if err != nil {
				c.Close(websocket.StatusInternalError, err.Error())
				return
			}
			//time.Sleep(1000*time.Mill isecond)
		}
	}()*/
	//c.Close(websocket.StatusNormalClosure, "")
}
