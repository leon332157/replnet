package main

// SNOW WAS HERE

import (
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/alecthomas/kong"
	koanfLib "github.com/knadh/koanf"
	koanfToml "github.com/knadh/koanf/parsers/toml"
	koanfBytes "github.com/knadh/koanf/providers/rawbytes"

	"github.com/leon332157/replish/client"
	"github.com/leon332157/replish/common"
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

var Command struct {
	Connect struct {
		RemoteURL string `arg:"" help:"remote url to connect to"`
		Port      uint16 `arg:"" help:"port to be forwarded" `
	} `cmd:"" help:"Connect to a repl" short:"c" optional:"" aliases:"c"`

	Serve struct {
		ListenPort uint16 `arg:"" help:"port to listen on" default:"0"`
		HttpPort   uint16 `arg:"" help:"port to listen on for http application" default:"0" optional:""`
	} `cmd:"" help:"Serve on a repl" optional:""`
	DefaultCommand struct{} `cmd:"" hidden:"" default:"1"`
	ConfigFile     string   `help:"Path to config file" default:".replit" short:"C"`
	LogLevel       string   `help:"Log Level (INFO,WARN,ERROR,DEBUG)" enum:"INFO, WARN, ERROR, DEBUG" default:"INFO" type:"enum"`
	Mode           string   `help:"Mode of operation, can be client or server" default:"client" short:"M" hidden:""`
}

func main() {
	ctx := kong.Parse(&Command, kong.Name("replish"), kong.Description("A websocket proxy for replit"))
	switch Command.LogLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	}
	log.Debugf("[Main] ctx.command: %s\n", ctx.Command())

	log.SetLevel(log.DebugLevel)
	switch ctx.Command() {
	case "connect": // assume we are connecting to a repl
		globalConfig.Mode = "client" // client
		globalConfig.RemoteURL = Command.Connect.RemoteURL
		globalConfig.RemoteAppPort = Command.Connect.Port
	case "serve": // assume we are serving a repl
		globalConfig.Mode = "server" // server
		globalConfig.ListenPort = Command.Serve.ListenPort
	case "default-command": // go to read config file
		ConfigFilePath := Command.ConfigFile
		log.Debugf("[Main] reading config file %s", ConfigFilePath)
		data, err := ioutil.ReadFile(ConfigFilePath)
		if err != nil {
			log.Fatalf("Failed to read config file: %v\n", err)
		}
		if err := loadConfigKoanf(data); err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
	}
	log.Debugf("[Main] running as %v", globalConfig.Mode)
	switch globalConfig.Mode {
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

//TODO: Maybe add both client and server field, then detect if running on replit to change mode or manually set modoe
// loadConfigKoanf loads the config file into koanf and checks for required configs
func loadConfigKoanf(content []byte) error {
	err := koanf.Load(koanfBytes.Provider(content), koanfToml.Parser())
	if err != nil {
		return err
	}

	// checks if replish field exist
	if koanf.Exists("replish") {
		koanf.Unmarshal("replish", &globalConfig)
		/*if err != nil {
			return fmt.Errorf("unmarshalling config failed: %v", err)
		}*/
	} else {
		return fmt.Errorf("replish field doesn't exist")
	}

	log.Debugf("[loadConfigKoanf] unmarshal %v\n", globalConfig)

	switch globalConfig.Mode {
	case "client":
		_, err := url.ParseRequestURI(globalConfig.RemoteURL) // check if url is valid
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
				return fmt.Errorf("local http port %v is invalid (1-65535)", appPort)
			}
			globalConfig.AppHttpPort = uint16(appPort)
		} else {
			log.Warn("local http port is not set")
		}
	default:
		return fmt.Errorf("mode %v is invalid", globalConfig.Mode)
		/*var prediction string
		if strings.ContainsAny(globalConfig.Mode, "svr") {
			prediction = "server"
		} else if strings.ContainsAny(globalConfig.Mode, "cli") {
			prediction = "client"
		} else {

		}
		return fmt.Errorf("mode %v is invalid, did you mean %v?", globalConfig.Mode, prediction)*/
	}
	log.Debugf("[loadConfigKoanf] Config:\n%s", koanf.Sprint())
	log.Debugf("[loadConfigKoanf] Effective Config:\n%v", globalConfig)
	return nil
}
