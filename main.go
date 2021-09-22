package main

// SNOW WAS HERE

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	koanfLib "github.com/knadh/koanf"
	koanfToml "github.com/knadh/koanf/parsers/toml"
	koanfBytes "github.com/knadh/koanf/providers/rawbytes"

	"github.com/leon332157/replish/client"
	"github.com/leon332157/replish/netstat"
	"github.com/leon332157/replish/server"

	log "github.com/sirupsen/logrus"
)

var (
	globalConfig ReplishConfig
	koanfConfig  = koanfLib.Conf{
		Delim:       ".",
		StrictMerge: true,
	}
	koanf = koanfLib.NewWithConf(koanfConfig)
)

const (
	ModeHelpString = "Mode of operation, can be client or server"
	UrlHelpString  = "URL of the repl (repl.co link)"
	ConfHelpString = "Path to config file"
)

type DotReplit struct {
	Run      string
	Language string
	OnBoot   string
	Packager map[string]interface{}
	Replish  ReplishConfig `koanf:"replish"`
}

type ReplishConfig struct {
	Mode           string `koanf:"mode"` // Mode of operation
	RemoteURL      string //`koanf:"remote-url"`     // The repl.co url to connect to
	LocalAppPort   uint16 //`koanf:"local-app-port"` // The port of your application
	RemoteAppPort  uint16 //`koanf:"remote-port"`    // The port of a remote application
	ListenPort     uint16 `koanf:"listen-port"` // The port replish listen on for WS connection
	configFilePath string
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetReportCaller(false)
	log.SetLevel(log.DebugLevel)
}

func startBasicHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello: %s", req.URL.Path)
	})
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func main() {
	parser := argparse.NewParser("replish", "Command line tool for replit")
	configFilePath := parser.String("C", "config", &argparse.Options{Help: ConfHelpString, Default: ".replit"})
	/*mode := parser.Selector("m", "mode", []string{"c", "client", "s", "server"}, &argparse.Options{Help: ModeHelpString, Default: "client"})
	replUrl := parser.String("c", "remote-url", &argparse.Options{Help: UrlHelpString, Default: "https://replit.com"})
	listenPort := parser.Int("p", "listen-port", &argparse.Options{Help: "Port to listen on", Default: 8080})
	*/
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		// fmt.Println(err)
		fmt.Print(parser.Usage(err))
		log.Exit(1)
	}
	globalConfig.configFilePath = *configFilePath
	if err := loadConfigKoanf(readConfigFile(globalConfig.configFilePath)); err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	// go startBasicHttp()
	// time.Sleep(1 * time.Second) // wait for server to come online
	// getPort()
	// log.Debugf("[Main] Got port: %v\n", port)
	go server.StartMain(7777, globalConfig.LocalAppPort)
	go client.ConnectWS(globalConfig.RemoteURL, globalConfig.RemoteAppPort, 10*time.Second)
	/*run, ok := dotreplit.Replish["run"].(string)
	if !ok {
		log.Warn("Reading 'run' field failed")
	}
	go client.ExecCommand(run)*/
	for {
		time.Sleep(1 * time.Second)
	}
}

func getPortAuto() uint16 {
	var port uint16
	addrs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		if s.Process == nil { // Process can be nil, discard it
			return false
		} else if strings.Contains(s.Process.Name, "System") {
			return false // Discard System process
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
		sel, err := strconv.ParseInt(inp, 10, 64)
		if err != nil {
			log.Panic(err)
		}
		if sel > int64(len(addrs)) { // Input is a port selection
			port = checkPort(sel)
		} else { // Input is list index
			temp := addrs[sel-1]
			port = temp.LocalAddr.Port
		}
	} else {
		port = addrs[0].LocalAddr.Port
	}
	return port
}

func checkPort(p int64) uint16 {
	if p > 65535 || p < 1 {
	
		log.Fatalf("port %v is out of range(1-65535)", p)
	}
	return uint16(p)
}

/*
func getPort() {
	log.Debug("Getting port")
	rawPort, ok := //dotreplit.Replish["port"] // Check if port exist
	if !ok {
		log.Warn("Port is missing, defaulting to auto")
		getPortAuto()
		return
	}
	intPort, ok := rawPort.(int64)
	if ok { // port is int
		port = checkPort(intPort)
		return
	}
	strPort, ok := rawPort.(string)
	if ok {
		// Port is string
		if strPort == "auto" {
			getPortAuto()
			return
		} else {
			intPort, err := strconv.Atoi(strPort)
			temp := int64(intPort)
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
*/
func readConfigFile(filepath string) []byte {
	log.Infof("reading config file %s", filepath)
	ioutil.ReadFile(filepath)
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading config file: %v\n", err)
	}
	return data
}

func loadConfigKoanf(content []byte) error {
	err := koanf.Load(koanfBytes.Provider(content), koanfToml.Parser())
	if err != nil {
		return err
	}
	if koanf.Exists("replish") {
		err = koanf.Unmarshal("replish", &globalConfig)
		if err != nil {
			return fmt.Errorf("unmarshalling replish failed: %v", err)
		}
	} else {
		return fmt.Errorf("replish field doesn't exist")
	}
	log.Debugln(koanf.Sprint())
	log.Debugln(globalConfig)

	//* check mode
	if len(globalConfig.Mode) < 4 {
		log.Warnln("mode is missing or invalid, defaulting to client")
		globalConfig.Mode = "client"
	}
	log.Infof("running as %v", globalConfig.Mode)
	if globalConfig.Mode == "client" {
		url, err := url.ParseRequestURI(globalConfig.RemoteURL)
		if err == nil {
			log.Debugln(url)
			// TODO: maybe check url speficis host and port
			//  Check that the remote url and remote app port are set
		} else {
			return fmt.Errorf("remote URL is not valid: %v", err)
		}
		if koanf.Exists("replish.remote-app-port") {
			remoteAppPort := koanf.Int64("replish.remote-app-port")
			globalConfig.RemoteAppPort = checkPort(remoteAppPort)
		} else {
			return fmt.Errorf("remote application port is unset")
		}
	} else if globalConfig.Mode == "server" { // Check that local app port is set
		if koanf.Exists("replish.listen-port") {
			listenPort := koanf.Int64("replish.listen-port")
			if listenPort > 65535 || listenPort < 0 {
				log.Warnln("listen port is invalid (0-65535), defaulting to 0")
				globalConfig.ListenPort = 0
			}
		} else {
			log.Warnln("listen port is unset, defaulting to 0")
			globalConfig.ListenPort = 0
		}

		if koanf.Exists("replish.local-app-port") {
			appPort := koanf.Int64("replish.app-port")
			globalConfig.LocalAppPort = checkPort(appPort)
		} else {
			return fmt.Errorf("local application port is not set")
		}
	} else {
		return fmt.Errorf("mode %v is invalid", globalConfig.Mode)
	}
	return nil
}
