package client

import (
	"time"

	"github.com/leon332157/replnet/common"
	log "github.com/sirupsen/logrus"
)

func StartMain(config *common.ReplnetConfig) {
	if len(config.Client.RemoteURL) == 0 {
		log.Fatalf("[Client Config] remote url len 0")
	}
	log.Debugf("[Client Config] remote url %v", config.Client.RemoteURL)
	if config.Client.RemotePort == 0 {
		log.Fatalf("[Client Config] remote app port is 0??")
	}
	log.Debugf("[Client Config] remote app port %v", config.Client.RemotePort)
	startWS(config.Client.RemoteURL, config.Client.RemotePort, config.Client.LocalPort, 10*time.Second)
}
