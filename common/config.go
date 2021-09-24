package common

type ReplishConfig struct {
	Mode           string `koanf:"mode"` // Mode of operation
	RemoteURL      string //`koanf:"remote-url"`     // The repl.co url to connect to
	LocalAppPort   uint16 //`koanf:"local-app-port"` // The port of your application
	RemoteAppPort  uint16 //`koanf:"remote-port"`    // The port of a remote application
	ListenPort     uint16 `koanf:"listen-port"` // The port replish listen on for WS connection
	ConfigFilePath string
}