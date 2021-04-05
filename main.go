package main

import (
	//"bufio"
	. "fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	fiber "github.com/gofiber/fiber/v2"
	server "github.com/leon332157/replish/server"
	toml "github.com/pelletier/go-toml"
)

var (
	dotreplit DotReplit
	port      uint16
)

type DotReplit struct {
	Run      string
	Language string
	Replish  map[string]interface{}
}

func main() {
	loadDotreplit()
	Println(dotreplit)
	go startFiber()
	time.Sleep(1 * time.Second) // wait for server to be created
	port = getPort()
	Printf("Got port: %v\n", port)
	go server.StartForwardServer(port)
	for {
		time.Sleep(1 * time.Second)
	}
}

func getPort() uint16 {
	var rawPort interface{}
	rawPort, ok := dotreplit.Replish["port"] // Check if port exist
	if !ok {
		Println("WARNING: Port is missing, defaulting to auto")
		rawPort = "auto"
	}
	_, ok = rawPort.(int64) // cheeck if port is int
	if ok {
		if rawPort.(int64) > 65535 || rawPort.(int64) < 1 {
			panic("Port out of range!")
		}
		return uint16(rawPort.(int64)) // port is int, return as uint16
	} else {
		// port is string
		rawPort, ok = dotreplit.Replish["port"].(string) // Check if port is string and exist
		if !ok {
			rawPort = "auto" // failed?? defaulting to auto
		}
		port, err := strconv.ParseUint(rawPort.(string), 10, 16)
		if err == nil {
			return uint16(port)
		} else {
			Printf("Error when converting port:%v, defaulting to auto\n", err)
			rawPort = "auto"
		}
	}

	if rawPort == "auto" {
		addrs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
			return net.IP.IsLoopback(s.LocalAddr.IP) && s.State == netstat.Listen
		})
		if err != nil {
			panic("Reading ports failed.")
		}
		if len(addrs) == 0 {
			panic("Looks like we aren't finding any open ports, are you listening on localhost (127.0.0.1)?")
		}
		for _, e := range addrs {
			if e.Process != nil {
				return e.LocalAddr.Port
			}
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

func loadDotreplit() {
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
	dotreplit = DotReplit{}
	err = toml.Unmarshal(contents, &dotreplit)
	if err != nil {
		panic(Sprintf("failed to unmarshal: %v\n", err))
	}
	if dotreplit.Replish == nil {
		panic("Replish field is empty! Check for typos in .replit!")
	}
}
