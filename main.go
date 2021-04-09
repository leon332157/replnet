package main

import (
	//"bufio"
	"fmt"
	"github.com/cakturk/go-netstat/netstat"
	fiber "github.com/gofiber/fiber/v2"
	server "github.com/leon332157/replish/server"
	toml "github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	dotreplit       DotReplit
	port            uint16
	hasReplishField bool = false
)

type DotReplit struct {
	Run      string
	Language string
	onBoot   string
	packager map[string]interface{}
	Replish  map[string]interface{}
}

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetReportCaller(false)
	log.SetLevel(log.DebugLevel)
	loadDotreplit()
	//go startHijack()
	go startFiber()
	time.Sleep(1*time.Second) // wait for server to come online
	getPort()
	log.Debugf("Got port: %v\n", port)
	go server.StartForwardServer(port)
	for {
		time.Sleep(1 * time.Second)
	}
}

func startHijack() {
	http.HandleFunc("/hijack", func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}
		conn, bufrw, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Don't forget to close the connection:
		defer conn.Close()
		bufrw.WriteString("Now we're speaking raw TCP. Say hi: ")
		bufrw.Flush()
		s, err := bufrw.ReadString('\n')
		if err != nil {
			log.Printf("error reading string: %v", err)
			return
		}
		fmt.Fprintf(bufrw, "You said: %q\nBye.\n", s)
		bufrw.Flush()
	})
	http.ListenAndServe(":8484", nil)
	log.Println("started")
}

func getPortAuto() {
	addrs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return net.IP.IsLoopback(s.LocalAddr.IP) && s.State == netstat.Listen
	})
	if err != nil {
		log.Fatalf("Reading ports failed:%v", err)
	}
	if len(addrs) == 0 {
		log.Fatalf("Looks like we aren't finding any open ports, are you listening on localhost (127.0.0.1)?")
	}
	for _, e := range addrs {
		if e.Process != nil {
			port = e.LocalAddr.Port
		}
	}
}

func getPort() {
	log.Debug("Getting port")
	rawPort, ok := dotreplit.Replish["port"] // Check if port exist
	if !ok {
		log.Warn("Port is missing, defaulting to auto")
		getPortAuto()
		return
	}
	intPort, ok := rawPort.(int64)
	if ok { // port is int
		if intPort > 65535 || intPort < 1 {
			log.Fatalf("port %v is out of range(1-65535)", rawPort)
		}
		port = uint16(intPort)
	}
	strPort, ok := rawPort.(string)
	if ok {
		// Port is string
		if strPort == "auto" {
			getPortAuto()
			return
		} else {
			temp, err := strconv.ParseUint(strPort, 10, 16)
			if err == nil {
				port = uint16(temp)
			} else {
				log.Errorf("Error when converting port: %v, defaulting to auto\n", err)
				getPortAuto()
				return
			}
		}
	}
}

func startFiber() {
	app := fiber.New(fiber.Config{DisableStartupMessage:false})

	app.Get("/*", func(c *fiber.Ctx) error {
		fmt.Println(c.Request())
		return c.SendString("haha")
	})

	go app.Listen("127.0.0.1:8383")
	log.Debug("fiber started")
}

func loadDotreplit() {
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
	dotreplit = DotReplit{}
	err = toml.Unmarshal(contents, &dotreplit)
	if err != nil {
		log.Panicf("failed to unmarshal: %v\n", err)
	}
	if dotreplit.Replish == nil {
		log.Warn("Replish field is empty or doesn't exist! Check for typos in .replit")
		// Write replish field maybe
	} else {
		hasReplishField = true
	}
}
