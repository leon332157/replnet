package server
import (
  "net/http"
  "log"
  "fmt"
)

func StartMain() {
  handler := http.HandlerFunc(handlerDav)
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", r.URL.Path)
})
	http.Handle("/__dav", handler)
	//http.FileServer(http.Dir("/home/runner/replish"))
	log.Fatal(http.ListenAndServe(":7878", nil))
}