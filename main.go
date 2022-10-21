package main

// SNOW WAS HERE

import (
	"fmt"
	"os"
	"net/url"

	"github.com/alecthomas/kong"
	koanfLib "github.com/knadh/koanf"
	koanfToml "github.com/knadh/koanf/parsers/toml"
	koanfRawBytes "github.com/knadh/koanf/providers/rawbytes"

	"github.com/leon332157/replnet/client"
	"github.com/leon332157/replnet/common"
	"github.com/leon332157/replnet/server"

	log "github.com/sirupsen/logrus"
)

var (
	globalConfig common.ReplnetConfig
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
	Run        string
	Entrypoint string
	OnBoot     string
	Compile    string
	//Packager map[string]interface{}
	Replnet common.ReplnetConfig `koanf:"replnet"`
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
	cmdCtx := kong.Parse(&Command, kong.Name("replish"), kong.Description("A websocket proxy for replit"))
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
	log.Debugf("[Main] ctx.command: %s\n", cmdCtx.Command())

	log.SetLevel(log.DebugLevel)
	switch cmdCtx.Command() { // read the command given
	case "connect": // the command given is to connect
		globalConfig.Mode = "client" // client
		globalConfig.Client.RemoteURL = Command.Connect.RemoteURL
		globalConfig.Client.RemotePort = Command.Connect.Port
	case "serve": // the command given is server
		globalConfig.Mode = "server" // set mode to server
		globalConfig.Server.ListenPort = Command.Serve.ListenPort
	case "default-command": // if none go to read config file
		ConfigFilePath := Command.ConfigFile
		log.Debugf("[Main] reading config file %s", ConfigFilePath)
		data, err := os.ReadFile(ConfigFilePath) // read the config file
		if err != nil {
			log.Fatalf("Failed to read config file: %v\n", err)
		}
		if err := loadConfigKoanf(data); err != nil { // load the config file
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

// TODO: Maybe add both client and server field, then detect if running on replit to change mode or manually set modoe
// loadConfigKoanf loads the config file into koanf and checks for required configs
func loadConfigKoanf(content []byte) error {
	err := koanf.Load(koanfRawBytes.Provider(content), koanfToml.Parser())
	if err != nil {
		return err
	}

	// checks if replish field exist
	if !koanf.Exists("replnet") {
		return fmt.Errorf("replnet field doesn't exist")
	}

	// checks if mode key exist``
	if !koanf.Exists("replnet.mode") {
		return fmt.Errorf("mode is required")
	}

	mode := koanf.MustString("replnet.mode")
	switch mode {
	case "client":
		globalConfig.Mode = "client"
		if !koanf.Exists("replnet.client") {
			return fmt.Errorf("replnet.client field doesn't exist")
		}

		if koanf.Exists("replnet.client.remote-url") {
			remoteURL := koanf.MustString("replnet.client.remote-url")
			parsedUrl, err := url.ParseRequestURI(remoteURL) // check if url is valid
			if err != nil {
				return fmt.Errorf("remote-url is invalid: %v", err)
			}

			log.Debugf("%+v\n", parsedUrl)
			if (parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https") &&
				(parsedUrl.Scheme != "ws" && parsedUrl.Scheme != "wss") {
				return fmt.Errorf("remote-url must be http, https, ws or wss")
			}
			if parsedUrl.Host == "" {
				return fmt.Errorf("remote-url is invalid: host is empty")
			}
			globalConfig.Client.RemoteURL = remoteURL
		} else {
			return fmt.Errorf("remote-url is required")
		}

		if koanf.Exists("replnet.client.remote-port") {
			remotePort := koanf.Int64("replnet.client.remote-port") // must use int64 to avoid overflow
			if remotePort > 65535 || remotePort < 1 {
				return fmt.Errorf("remote-port %v is invalid (1-65535)", remotePort)
			}
			globalConfig.Client.RemotePort = uint16(remotePort)
		} else {
			return fmt.Errorf("remote-port is required")
		}

		if koanf.Exists("replnet.client.local-port") {
			localPort := koanf.Int64("replnet.client.local-port")

			if localPort > 65535 || localPort < 1 {
				return fmt.Errorf("local-port %v is invalid (1-65535)", localPort)
			}
			globalConfig.Client.LocalPort = uint16(localPort)
		} else {
			globalConfig.Client.LocalPort = globalConfig.Client.RemotePort // default to remote app port if not set
			log.Warnf("local-port is not set, defaulting to remote-port: %v", globalConfig.Client.RemotePort)
		}

	case "server":
		globalConfig.Mode = "server"
		/*if !koanf.Exists("replnet.server") {
			return fmt.Errorf("replnet.server field doesn't exist")
		}*/

		if koanf.Exists("replnet.server.listen-port") {
			listenPort := koanf.Int64("replnet.server.listen-port")
			// Check if listen port is in valid range
			if listenPort > 65535 || listenPort < 0 {
				log.Warnln("listen port is invalid (0-65535), defaulting to 0")
				listenPort = 0
			}
			globalConfig.Server.ListenPort = uint16(listenPort)
		} else {
			log.Warnln("listen port is unset, defaulting to 0")
			globalConfig.Server.ListenPort = 0
		}

		if koanf.Exists("replnet.server.reverse-proxy-port") {
			reverseProxyPort := koanf.Int64("replnet.server.reverse-proxy-port")
			if reverseProxyPort > 65535 || reverseProxyPort < 0 {
				return fmt.Errorf("reverse-proxy-port %v is invalid (0-65535)", reverseProxyPort)
			}
			globalConfig.Server.ReverseProxyPort = uint16(reverseProxyPort)
		} else {
			log.Warn("reverse proxy port is not set, no reverse proxy will be used")
			globalConfig.Server.ReverseProxyPort = 0
		}

	default:
		return fmt.Errorf("mode %v is invalid", mode)
		/*var prediction string
		if strings.ContainsAny(globalConfig.Mode, "svr") {
			prediction = "server"
		} else if strings.ContainsAny(globalConfig.Mode, "cli") {
			prediction = "client"
		} else {

		}
		return fmt.Errorf("mode %v is invalid, did you mean %v?", globalConfig.Mode, prediction)*/
	}
	log.Debugf("[loadConfigKoanf] .replit:\n%s", koanf.Sprint())
	log.Debugf("[loadConfigKoanf] Loaded Global Config: %+v\n", globalConfig)
	return nil
}
