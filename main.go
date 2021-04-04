package main

import (
	//"bufio"
	"fmt"
	"html"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
	"github.com/cakturk/go-netstat/netstat"
	fiber "github.com/gofiber/fiber/v2"
	server "github.com/leon332157/replish/server"
	toml "github.com/pelletier/go-toml"
)

type DotReplit struct {
	Run      string
	Language string
	Replish  map[string]string
}

func main() {
	cfg := loadDotreplit()
	fmt.Println(cfg)
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
	http.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
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

func loadDotreplit() DotReplit {
	slug, ok := os.LookupEnv("REPL_SLUG")
	var path string
	if ok {
		path = fmt.Sprintf("/home/runner/%v/.replit", slug)
	} else {
		path = ".replit"
	}
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		contents = make([]byte, 0)
	}
	dotreplit := DotReplit{}
	err = toml.Unmarshal(contents, &dotreplit)
	if err != nil {
		fmt.Println(err)
	}
	return dotreplit
}
