package main

// SNOW WAS HERE

import (
	_ "bufio"
	"fmt"
	"io/ioutil"
	_ "net"
	"net/http"
	"net/url"
	"os"
	_ "strconv"
	"strings"

	"github.com/akamensky/argparse"
	koanfLib "github.com/knadh/koanf"
	koanfToml "github.com/knadh/koanf/parsers/toml"
	koanfBytes "github.com/knadh/koanf/providers/rawbytes"

	"github.com/leon332157/replish/client"
	"github.com/leon332157/replish/common"
	_ "github.com/leon332157/replish/netstat"
	"github.com/leon332157/replish/server"

	log "github.com/sirupsen/logrus"
)

var (
	globalConfig common.ReplishConfig
	koanfConfig  = koanfLib.Conf{
		Delim:       ".",
		StrictMerge: true,
	}
	koanf = koanfLib.NewWithConf(koanfConfig)
)

const (
	ModeHelpString     = "Mode of operation, can be client or server"
	UrlHelpString      = "URL of the repl (repl.co link)"
	ConfHelpString     = "Path to config file"
	LogLevelHelpString = "Level for logging"
)

type DotReplit struct {
	Run      string
	Language string
	OnBoot   string
	Packager map[string]interface{}
	Replish  common.ReplishConfig `koanf:"replish"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
	log.SetReportCaller(false)
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
	logLevel := parser.Selector("", "log-level", []string{"INFO", "WARN", "ERROR", "DEBUG"}, &argparse.Options{Default: "INFO"})
	serverFlag := parser.Flag("", "server", nil)
	/*mode := parser.Selector(
		"m",
		"mode",
		[]string{"c", "client", "s", "server"},
		&argparse.Options{Help: ModeHelpString, Default: "client"},
	)
	replUrl := parser.String("c", "remote-url", &argparse.Options{Help: UrlHelpString, Default: nil})
	listenPort := parser.Int("p", "listen-port", &argparse.Options{Help: "Port to listen on", Default: 8080})
	server.UNUSED(mode, replUrl, listenPort)
	*/

	// Parse input
	if err := parser.Parse(os.Args); err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		// fmt.Println(err)
		fmt.Print(parser.Usage(err))
		log.Exit(1)
	}
	switch *logLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	}
	globalConfig.ConfigFilePath = *configFilePath
	var content []byte
	if *serverFlag {
		content = []byte("[replish]\nmode = 'server'\nlocal-http-port=7777\nlisten-port = 9999")
	} else {
		content = readConfigFile(globalConfig.ConfigFilePath)
	}
	if err := loadConfigKoanf(content); err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	log.Infof("[Main] running as %v", globalConfig.Mode)
	switch strings.ToLower(globalConfig.Mode) {
	case "client":
		client.StartMain(&globalConfig)
	case "server":
		server.StartMain(&globalConfig)
	}
}

/*
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

// readConfigFile reads the config file and returns the config as a bytes
func readConfigFile(filepath string) []byte {
	log.Infof("[Main] reading config file %s", filepath)
	ioutil.ReadFile(filepath)
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading config file: %v\n", err)
	}
	return data
}

//TODO: Maybe add both client and server field, then detect if running on replit to change mode or manually set modoe 
// loadConfigKoanf loads the config file into koanf and checks for required configs
func loadConfigKoanf(content []byte) error {
	err := koanf.Load(koanfBytes.Provider(content), koanfToml.Parser())
	if err != nil {
		return err
	}

	// checks if replish field exist
	if koanf.Exists("replish") {
		err = koanf.Unmarshal("replish", &globalConfig)
		if err != nil {
			return fmt.Errorf("unmarshalling config failed: %v", err)
		}
	} else {
		return fmt.Errorf("replish field doesn't exist")
	}
	log.Debugln(koanf.Sprint())
	log.Debugln(globalConfig)

	switch strings.ToLower(globalConfig.Mode) {
	case "client":
		_, err := url.ParseRequestURI(globalConfig.RemoteURL)
		if err != nil {
			return fmt.Errorf("remote URL is invalid: %v", err)
		}
		if koanf.Exists("replish.remote-app-port") {
			remoteAppPort := koanf.Int64("replish.remote-app-port")
			if remoteAppPort > 65535 || remoteAppPort < 1 {
				return fmt.Errorf("remote app port %v is invalid (1-65535)", remoteAppPort)
			}
			globalConfig.RemoteAppPort = uint16(remoteAppPort)
		} else {
			return fmt.Errorf("remote application port is unset")
		}
		if koanf.Exists("replish.local-app-port") {
			localAppPort := koanf.Int64("replish.local-app-port")
			if localAppPort > 65535 || localAppPort < 1 {
				return fmt.Errorf("local app port %v is invalid (1-65535)", localAppPort)
			}
			globalConfig.LocalAppPort = uint16(localAppPort)
		} else {
			globalConfig.LocalAppPort = globalConfig.RemoteAppPort
		}
	case "server":
		if koanf.Exists("replish.listen-port") {
			listenPort := koanf.Int64("replish.listen-port")
			// Check if listen port is in valid range
			if listenPort > 65535 || listenPort < 0 {
				log.Warnln("listen port is invalid (0-65535), defaulting to 0")
				globalConfig.ListenPort = 0
			}
		} else {
			// Default to 0 if listen port doesn't exist
			log.Warnln("listen port is unset, defaulting to 0")
		}

		if koanf.Exists("replish.local-http-port") {
			appPort := koanf.Int64("replish.local-http-port") // read int64 because fool proof
			if appPort > 65535 || appPort < 1 {
				return fmt.Errorf("local http port is invalid (1-65535)")
			}
			globalConfig.AppHttpPort = uint16(appPort)
		} else {
			log.Warn("local http port is not set")
		}
	default:
		var prediction string
		if strings.ContainsAny(globalConfig.Mode, "svr") {
			prediction = "server"
		} else if strings.ContainsAny(globalConfig.Mode, "cli") {
			prediction = "client"
		} else {
			return fmt.Errorf("mode %v is invalid", globalConfig.Mode)
		}
		return fmt.Errorf("mode %v is invalid, did you mean %v?", globalConfig.Mode, prediction)
	}

	return nil
}
