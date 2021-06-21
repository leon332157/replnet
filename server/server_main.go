package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

//TODO: Handler for __dav, *.git, __ws, __ssh and wildcard (reverse proxy)
func StartMain(listenPort, forwardPort uint16) {
	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	})

	http.HandleFunc("/__dav", handlerDav)
	*/
	//http.FileServer(http.Dir("/home/runner/replish"))
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%v", listenPort))

	if err != nil {
		log.Panicf("[Server Main] %s\n", err)
	}
	log.Infof("[Server Main] Listening on %v", listenPort)
	//p := &ReverseProxy{port: port}
	http.Serve(listener, &ReplishRouter{port: forwardPort})
	/*go func() {=
		err := http.Serve(listener, p)

		if err != nil {
			log.Panicf("[Server Main] %s\n", err)
		}
	}()

	if err != nil {
		log.Panicf("[Server Main] %s\n", err)
	}
	log.Debug("[Server Main] reverse proxy started")*/
}

type ReplishRouter struct {
	port uint16
}

func (s *ReplishRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/__dav") {
		log.Debug("[Server Router] Match /__dav, passing to webdav")
		handlerDav(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/__ws") {
		log.Debug("[Server Router] Matching /__ws, passing to websocket")
		handleWS(w,r)
	} else {
		localUrl, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%v", s.port))
		if err != nil {
			log.Fatalf("[Server Router] Formatting url failed!")
		}
		proxy := httputil.NewSingleHostReverseProxy(localUrl)
		proxy.ServeHTTP(w, r)
	}
}
