package main

import (
	"time"

	"github.com/blend/go-sdk/logger"
)

func main() {
	log := logger.MustNew(logger.OptConfigFromEnv())
	tick := time.Tick(time.Second)
	for {
		select {
		case <-tick:
			log.Infof("it's %s", time.Now().Format(time.RFC3339))
		}
	}
}
