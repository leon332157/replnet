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
func StartMain(port uint16) {
	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	})
  
	http.HandleFunc("/__dav", handlerDav)
  */
	//http.FileServer(http.Dir("/home/runner/replish"))
	listener, err := net.Listen("tcp4", ":0")

	if err != nil {
		log.Panicf("[Server Main] %s\n", err)
	}

	//p := &ReverseProxy{port: port}
  http.Serve(listener,&ReplishRouter{port:port})
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
  if strings.HasPrefix(r.URL.Path,"/__dav") {
    log.Debug("[Server Router] Match /__dav, passing to webdav")
    handlerDav(w,r)
  } else {
    fmt.Fprintf(w, "Hello, %q", r.URL.Path)
    localUrl, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%v", p.port))
	if err != nil {
    log.FatalFn("[Server Router] Formatting url failed!")
	}
    proxy := httputil.NewSingleHostReverseProxy(localUrl)
    proxy.ServeHTTP(w,r)
  }
}

