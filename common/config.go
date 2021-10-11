package common

import (
	"fmt"
	_  "github.com/sirupsen/logrus"
	_ "net/url"
)

type ReplishConfig struct {
	Mode           string `koanf:"mode"`            // Mode of operation
	RemoteURL      string `koanf:"remote-url"`      // The repl.co url to connect to
	LocalHttpPort  uint16 `koanf:"local-app-port"`  // The port of your http application
	RemoteAppPort  uint16 `koanf:"remote-app-port"` // The port of a remote application
	ListenPort     uint16 `koanf:"listen-port"`     // The port replish listen on for WS connection
	ConfigFilePath string
}

// checkPort takes a int value, checkes 1-65535 and converts to uint16 if no error
func checkPort(p int64) (uint16, error) {
	if p > 65535 || p < 1 {
		return 0, fmt.Errorf("port %v is out of range(1-65535)", p)
	}
	return uint16(p), nil
}
/*
// check config checks for required fields when in different modes
func (c *ReplishConfig) checkConfig() error {
	switch c.Mode {
	case "client":
		url, err := url.ParseRequestURI(c.RemoteURL)
		if err == nil {
			log.Debugf("RemoteURL: %v", url)
			// TODO: maybe check url speficis host and port
		} else {
			return fmt.Errorf("remote URL is not valid: %v", err)
		}
		// check remote app port
		if c.RemoteAppPort == 0 {
			return fmt.Errorf("remote application port is unset")
		}
	case "server": // Check that local app port is set
		if c.LocalAppPort == 0 {
			return fmt.Errorf("local application port is not set")
		}
	default:
		return fmt.Errorf("mode %v is invalid", c.Mode)
	}
	return nil
}
*/