package server

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/leon332157/replnet/common"
	log "github.com/sirupsen/logrus"
)

// TODO: Handler for __dav, *.git, __ws, __ssh and wildcard (reverse proxy)
func StartMain(config *common.ReplnetConfig) {
	// check server configs
	if config.Server.ReverseProxyPort == 0 {
		log.Warnln("[Server Config] app http port is 0, running without reverse proxy")
	}
	//http.FileServer(http.Dir("/home/runner/replish"))
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%v", config.Server.ListenPort))
	if err != nil {
		log.Panicf("[Server Main] %s\n", err)
	}
	log.Infof("[Server Main] Listening on %v", listener.Addr().String())
	// p := &ReverseProxy{port: port}
	http.Serve(listener, &ReplishRouter{config: config})
}

type ReplishRouter struct {
	config *common.ReplnetConfig
}

func (s *ReplishRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/__webdav") {
		log.Debug("[Server Router] Match /__webdav, passing to webdav")
		handlerDav(w, r)
	} else if strings.HasPrefix(path, "/__ws") {
		log.Debug("[Server Router] Matching /__ws, passing to websocket")
		handleWS(w, r)
	} else if strings.HasPrefix(path, "/__ping") {
		w.Write([]byte("pong"))
	} else {
        if s.config.Server.ReverseProxyPort == 0 {
			log.Debug("[Server Router] Reverse proxy port is 0, returning 204")
            w.WriteHeader(http.StatusNoContent)
			return
		}
        //TODOL only return 204 for now
        log.Debugf("[Server Router] Matching %s, passing to reverseProxy", path)
		localUrl, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%v", s.config.Server.ReverseProxyPort))
		if err != nil {
			log.Fatalf("[Server Router] Formatting url failed!")
		}
		proxy := httputil.NewSingleHostReverseProxy(localUrl)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Errorf("[Server Reverse Proxy] %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		proxy.ServeHTTP(w, r)
	}
}
