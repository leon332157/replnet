package main

import (
	"bufio"
	"fmt"
	"github.com/leon332157/replish/netstat"
	fiber "github.com/gofiber/fiber/v2"
	server "github.com/leon332157/replish/server"
	toml "github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	time.Sleep(1 * time.Second) // wait for server to come online
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
		if s.Process == nil { // Process can be nil, discard it
			return false
		}
		return net.IP.IsLoopback(s.LocalAddr.IP) && s.State == netstat.Listen
	})
	if err != nil {
		log.Fatalf("Reading ports failed:%v", err)
	}
	if len(addrs) == 0 {
		log.Fatalf("Looks like we aren't finding any open ports, are you listening on localhost (127.0.0.1)?")
	} else if len(addrs) > 1 {
		fmt.Printf("Multiple ports detected: %v\n", len(addrs))
		for index, sock := range addrs {
			if sock.Process != nil {
				fmt.Printf("%v. %v %v\n", index+1, sock.Process, sock.LocalAddr.Port)
			}
		}
		fmt.Print("Choose port/process: ")
		reader := bufio.NewReader(os.Stdin)
		inp, err := reader.ReadString('\n')
		inp = strings.TrimSuffix(inp, "\r\n")
		if err != nil {
			log.Panic(err)
		}
		sel, err := strconv.Atoi(inp)
		if err != nil {
			log.Panic(err)
		}
		if sel > len(addrs) { // Input is a port selection
			port = checkPort(sel)
		} else { // Input is list index
			temp := addrs[sel-1]
			port = temp.LocalAddr.Port
		}
	} else {
		port = addrs[0].LocalAddr.Port
	}
}
func checkPort(p int) uint16 {
	if p > 65535 || p < 1 {
		log.Fatalf("port %v is out of range(1-65535)", p)
	}
	return uint16(p)
}
func getPort() {
	log.Debug("Getting port")
	rawPort, ok := dotreplit.Replish["port"] // Check if port exist
	if !ok {
		log.Warn("Port is missing, defaulting to auto")
		getPortAuto()
		return
	}
	intPort, ok := rawPort.(int)
	if ok { // port is int
		port = checkPort(intPort)
	}
	strPort, ok := rawPort.(string)
	if ok {
		// Port is string
		if strPort == "auto" {
			getPortAuto()
			return
		} else {
			temp, err := strconv.Atoi(strPort)
			if err == nil {
				port = checkPort(temp)
			} else {
				log.Errorf("Error when converting port: %v, defaulting to auto\n", err)
				getPortAuto()
				return
			}
		}
	}
}

func startFiber() {
	app := fiber.New(fiber.Config{DisableStartupMessage: false})

	app.Get("/*", func(c *fiber.Ctx) error {
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
