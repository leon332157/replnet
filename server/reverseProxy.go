package server

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func StartReverseProxy() {
	listener, err := net.Listen("tcp4", ":8484")

	if err != nil {
		panic(err)
	}

	p := &proxy{}

	go func() {
		err := http.Serve(listener, p)

		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}
	fmt.Println("reverse proxy started")
}

type proxy struct{}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	str, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%v%s", 7373, r.URL))
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(str)
	proxy.ServeHTTP(w, r)
}
