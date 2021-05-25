package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//TODO: Handler for __dav, *.git, __ws, __ssh and wildcard (reverse proxy)
func StartMain(port uint16) {
	handler := http.HandlerFunc(handlerDav)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	})
	http.Handle("/__dav", handler)
	//http.FileServer(http.Dir("/home/runner/replish"))
	listener, err := net.Listen("tcp4", ":8484")

	if err != nil {
		log.Panicf("[Server Main] %s\n", err)
	}

	p := &ReverseProxy{port: port}

	go func() {
		err := http.Serve(listener, p)

		if err != nil {
			log.Panicf("[Server Main] %s\n", err)
		}
	}()

	if err != nil {
		log.Panicf("[Server Main] %s\n", err)
	}
	log.Debug("[Server Main] reverse proxy started")
}

type ReverseProxy struct {
	port uint16
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	str, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%v", p.port))
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(str)
	proxy.ServeHTTP(w, r)
}
