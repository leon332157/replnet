package main

import (
	"bufio"
	"fmt"
	"html"
	"net/http"
	"os"

	"github.com/cakturk/go-netstat/netstat"
	fiber "github.com/gofiber/fiber/v2"
	server "github.com/leon332157/replish/server"
)

func main() {
	go server.StartForwardServer()
	readOpenTCP()
	readReplConfig()
	startHttp()
}

func readOpenTCP() error {
	addrs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	})
	if err != nil {
		return err
	}
	for _, e := range addrs {
		fmt.Printf("%v\n", e)
	}
	return nil
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

func readReplConfig() {
	path := fmt.Sprintf("/home/runner/%v/.replit", os.Getenv("REPL_SLUG"))
	//path := "main.go"
	var lines []string
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		f.Close()
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	fmt.Println(lines)
}
