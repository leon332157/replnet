package main

//SNOW WAS HERE

import (
	"bufio"
	"fmt"

	"github.com/akamensky/argparse"
	koanfLib "github.com/knadh/koanf"
	koanfToml "github.com/knadh/koanf/parsers/toml"
	koanfBytes "github.com/knadh/koanf/providers/rawbytes"

	"github.com/leon332157/replish/client"
	"github.com/leon332157/replish/netstat"

	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	port         uint16
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
	Mode           string `koanf:"mode"`        // Mode of operation
	RemoteURL      string `koanf:"remote-url"`  // The repl.co url to connect to
	LocalPort      uint16 `koanf:"local-port"`  // The port of your application
	RemotePort     uint16 `koanf:"remote-port"` // The port of a remote application
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
	/*configFile := parser.String("C", "config", os.O_RDONLY, 0777, &argparse.Options{Help: ConfHelpString, Default: ".replit"})
	mode := parser.Selector("m", "mode", []string{"c", "client", "s", "server"}, &argparse.Options{Help: ModeHelpString, Default: "client"})
	replUrl := parser.String("c", "remote-url", &argparse.Options{Help: UrlHelpString, Default: "https://replit.com"})
	listenPort := parser.Int("p", "listen-port", &argparse.Options{Help: "Port to listen on", Default: 8080})
	*/
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		//fmt.Println(err)
		fmt.Print(parser.Usage(err))
		log.Exit(1)
	}
	globalConfig.configFilePath = *configFilePath
	if err := loadConfigKoanf(readConfigFile(globalConfig.configFilePath)); err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	//go startBasicHttp()
	//time.Sleep(1 * time.Second) // wait for server to come online
	//getPort()
	port = 8080
	//log.Debugf("[Main] Got port: %v\n", port)
	//go server.StartMain(7777, port)
	go client.ConnectWS(globalConfig.RemoteURL, globalConfig.RemotePort, 10*time.Second)
	/*run, ok := dotreplit.Replish["run"].(string)
	if !ok {
		log.Warn("Reading 'run' field failed")
	}
	go client.ExecCommand(run)*/
	for {
		time.Sleep(1 * time.Second)
	}
}

func getPortAuto() {
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
	// check config
	if globalConfig.Mode == "" {
		log.Warnln("mode is missing, defaulting to client")
		globalConfig.Mode = "client"
	}
	if globalConfig.Mode == "client" {
		if _, err := url.ParseRequestURI(globalConfig.RemoteURL); err != nil {
			//  Check that the remote url and remote app port are set
			return fmt.Errorf("remote URL is not valid: %v", err)
		}
		if !koanf.Exists("replish.remote-port") {
			return fmt.Errorf("remote port is unset")

		}
		remotePort := koanf.Int64("replish.remote-port")
		if remotePort > 65535 || remotePort < 1 {
			return fmt.Errorf("remote port is invalid (1-65535)")
		}
	}
	if globalConfig.Mode == "server" { // Check that local app port is set
		if !koanf.Exists("replish.listen-port") {
			log.Warnln("listen port is unset, defaulting to 0")
			globalConfig.ListenPort = 0
		}
		listenPort := koanf.Int64("replish.listen-port")
		if listenPort > 65535 || listenPort < 0 {
			log.Warnln("listen port is invalid (0-65535), defaulting to 0")
			globalConfig.ListenPort = 0
		}
		if koanf.Exists("replish.local-port") {
		} else {
			return fmt.Errorf("local port is not set")
		}
		localPort := koanf.Int64("replish.local-port")
		if localPort == 0 {
			//TODO: Attempt auto detect local port
		} else {
			if localPort > 65535 || localPort < 1 {
				return fmt.Errorf("local port is invalid (1-65535)")
			}
		}
	}
	return nil
}
