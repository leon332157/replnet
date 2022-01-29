package common

type ReplishConfig struct {
	Mode           string `koanf:"mode"`            // Mode of operation
	RemoteURL      string `koanf:"remote-url"`      // Client: The repl.co url to connect to
	AppHttpPort    uint16 `koanf:"app-http-port"`   // Server: OPTINAL: The port of your http application
	RemoteAppPort  uint16 `koanf:"remote-app-port"` // Client: The port of the remote(tcp) application
	LocalAppPort   uint16 `koanf:"local-app-port"`  // Client OPTINAL: The port of client listener to proxy to remote app
	ListenPort     uint16 `koanf:"listen-port"`     // Server OPTINAL: The port replish server listen on
}