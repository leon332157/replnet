package server

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"golang.org/x/net/webdav"
	log "github.com/sirupsen/logrus"
)

func handleBasicAuth(username string, password string) bool {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	USERNAME, ok := os.LookupEnv("REPLISH_USERNAME")
	if !ok {
		log.Fatal("Looking up username failed, defaulting to test")
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
		http.Error(w, "Not authorized", http.StatusUnauthorized)
	}
	log.Printf("username: %s\npassword: %s\n", username, password)
	if !handleBasicAuth(username, password) {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
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
	//TODO: Maybe add directory listing

	srv := &webdav.Handler{
		Prefix:     "/__dav",
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
	//UNUSED(srv)
	srv.ServeHTTP(w, r)
}
