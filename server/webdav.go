package server

import (
	_ "github.com/leon332157/replish/server/webdav"
  "log"
  "net/http"
  "golang.org/x/net/webdav"
)
func StartWebdav() {
  srv := &webdav.Handler{
		FileSystem: webdav.Dir("/home/runner/replish"),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				log.Printf("WEBDAV [%s]: %s \n", r.Method, r.URL)
			}
		},
	}
	http.Handle("/", srv)
  //http.FileServer(http.Dir("/home/runner/replish"))
	log.Fatal(http.ListenAndServe(":8080", nil))
}