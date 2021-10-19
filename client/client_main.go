package client

import (
	"time"

	"github.com/leon332157/replish/common"
	log "github.com/sirupsen/logrus"
)

func StartMain(config *common.ReplishConfig) {
	if len(config.RemoteURL) == 0 {
		log.Fatalf("[Client Config] remote url len 0")
	}
	log.Debugf("[Client Config] remote url %v", config.RemoteURL)
	if config.RemoteAppPort == 0 {
		log.Fatalf("[Client Config] remote app port is 0??")
	}
	log.Debugf("[Client Config] remote app port %v", config.RemoteAppPort)
	startWS(config.RemoteURL, config.RemoteAppPort, config.LocalAppPort, 10*time.Second)
}
