package server

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/webdav"
	"net/http"
	"os"
	"strings"
)
var (
    string ROOT_PATH = nil
    webdav.Handler webDavHandler = &webdav.Handler{
		Prefix:     "/__webdav",
		FileSystem: webdav.Dir(ROOT_PATH),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Debugf("[Webdav Handler] %s %s: %s, ERROR: %s\n", r.UserAgent(), r.Method, r.URL, err)
			} else {
				log.Debugf("[Webdav Handler] %s %s: %s \n", r.UserAgent(), r.Method, r.URL)
			}
		},
	}
)

func handleBasicAuth(username string, password string) bool {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	USERNAME, ok := os.LookupEnv("REPLISH_USER")
	if !ok {
		log.Error("Looking up username failed")
		return true
	}
	PASSWORD, ok := os.LookupEnv("REPLISH_PW")
	if !ok {
		log.Error("Looking up password failed")
		return true
	}
	return username == USERNAME && password == PASSWORD

}
func handlerDav(w http.ResponseWriter, r *http.Request) {
    log.Debugf("[Webdav Handler] URL: %#v",r.URL)
	if !strings.Contains(r.URL.Path, ".git") {
		username, password, ok := r.BasicAuth()
		if !ok {
			log.Debug("[Webdav Handler] BasicAuth error")
			w.Header().Set("WWW-Authenticate", `Basic realm="Enter username and password"`)
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		//log.Printf("username: %s\npassword: %s\n", username, password)
		if !handleBasicAuth(username, password) {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		} else {
			log.Debug("[Webdav Handler] auth passed")
		}
	}
	var ROOT_PATH string
	REPL_SLUG, ok := os.LookupEnv("REPL_SLUG")
	if ok {
		ROOT_PATH = fmt.Sprintf("/home/runner/%s", REPL_SLUG)
	} else {
		ROOT_PATH, _ = os.Getwd()
	}
	//TODO: Maybe add directory listing
    
	webDavHandler.ServeHTTP(w, r)
}
