package server

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func StartReverseProxy() {
	listener, err := net.Listen("tcp", ":8484")

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

}

type proxy struct{}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parsed, err := url.Parse(fmt.Sprintf("http://127.0.0.1:7373%s", r.URL))
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(parsed)
	proxy.ServeHTTP(w, r)
}
