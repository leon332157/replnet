package main

import (
	"fmt"
	"html"
	"net/http"

	fiber "github.com/gofiber/fiber/v2"
	server "github.com/leon332157/replish/server"
)

func main() {
	go server.StartForwardServer()
	startHttp()
}

func startHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})
	err := http.ListenAndServe("127.0.0.1:8181", nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func startFiber() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Listen(":3000")
}
