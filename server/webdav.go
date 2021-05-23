package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/leon332157/replish/server/webdav"
	"golang.org/x/net/webdav"
)

func handleBasicAuth(username string, password string) bool {
	USERNAME, ok := os.LookupEnv("REPLISH_USERNAME")
	if !ok {
		log.Fatal("Looking up username failed")
		return false
	}
	PASSWORD, ok := os.LookupEnv("REPLISH_PASSWORD")
	if !ok {
		log.Fatal("Looking up password failed")
		return false
	}
	return username == USERNAME && password == PASSWORD

}
func handlerDav(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic`)
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Not authorized", 401)
	}
	log.Printf("username: %s\npassword: %s\n", username, password)
	if !handleBasicAuth(username, password) {
		http.Error(w, "Not authorized", 401)
	} else {
		fmt.Println("passed")
	}
	var ROOT_PATH string
	REPL_SLUG, ok := os.LookupEnv("REPL_SLUG")
	if ok {
		ROOT_PATH = fmt.Sprintf("/home/runner/%s", REPL_SLUG)
	} else {
		ROOT_PATH, _ = os.Getwd()
	}

	srv := &webdav.Handler{
		FileSystem: webdav.Dir(ROOT_PATH),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				log.Printf("WEBDAV [%s]: %s \n", r.Method, r.URL)
			}
		},
	}
	UNUSED(srv)
	//srv.ServeHTTP(w, r)
}
func StartWebdav() {
	handler := http.HandlerFunc(handlerDav)
	http.Handle("/", handler)
	//http.FileServer(http.Dir("/home/runner/replish"))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
