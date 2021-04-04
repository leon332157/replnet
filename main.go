package main

import (
	//"bufio"
	. "fmt"
	"github.com/cakturk/go-netstat/netstat"
	fiber "github.com/gofiber/fiber/v2"
	server "github.com/leon332157/replish/server"
	toml "github.com/pelletier/go-toml"
	"io/ioutil"
	"net"
	"os"
	"time"
)

type DotReplit struct {
	Run      string
	Language string
	Replish  map[string]string
}

func main() {
	cfg := loadDotreplit()
	Println(cfg)
	go startFiber()
	time.Sleep(1 * time.Second) // wait for server to be created
	port := readOpenTCP()
	Printf("Got port: %v\n", port)
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
		Println("Looks like we aren't finding any open ports, are you listening on localhost (127.0.0.1)?")
	}
	for _, e := range addrs {
		if e.Process != nil {
			//fmt.Printf("%v\n", e)
			return e.LocalAddr.Port
		}
	}
	return 0
}
func startFiber() {
	app := fiber.New()

	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendString("")
	})

	app.Listen("127.0.0.1:8383")
}

func loadDotreplit() DotReplit {
	slug, ok := os.LookupEnv("REPL_SLUG")
	var path string
	if ok {
		path = Sprintf("/home/runner/%v/.replit", slug)
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
		Println(err)
	}
	return dotreplit
}
