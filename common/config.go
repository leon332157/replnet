package common

import (
// semver "github.com/blang/semver/v4"
)

const VERSION = "0.1.0a"

type ReplnetConfig struct {
	Mode   string `koanf:"mode"` // Mode of operation, can be client or server
	Client ReplnetClientConfig
	Server ReplnetServerConfig
}
type ReplnetClientConfig struct {
	RemoteURL  string `koanf:"remote-url"`
	RemotePort uint16 `koanf:"remote-port"`
	LocalPort  uint16 `koanf:"local-port"`
}
type ReplnetServerConfig struct {
	ReverseProxyPort uint16 `koanf:"reverse-proxy-port"`
	ListenPort       uint16 `koanf:"listen-port"`
}
