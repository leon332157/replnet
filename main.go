package main

import (
	"bufio"
	"fmt"
	"html"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	fiber "github.com/gofiber/fiber/v2"
	server "github.com/leon332157/replish/server"
)

func main() {
	//readReplConfig()
	go startHttp()
	time.Sleep(1 * time.Second) // wait for server to be created
	port := readOpenTCP()
	fmt.Printf("Got port: %v\n", port)
	go server.StartForwardServer(port)
	for {

		time.Sleep(1 * time.Second)
	}
}

func readOpenTCP() uint16 {
	addrs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return net.IP.IsLoopback(s.LocalAddr.IP) && s.State == netstat.Listen
	})
	if err != nil {
		return 0
	}
	if len(addrs) == 0 {
		fmt.Println("Looks like we aren't finding any open ports, are you listening on localhost (127.0.0.1)?")
	}
	for _, e := range addrs {
		if e.Process != nil {
			//fmt.Printf("%v\n", e)
			return e.LocalAddr.Port
		}
	}
	return 0
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
	fmt.Println("http started")
}

func startFiber() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Listen("127.0.0.1:8383")
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
