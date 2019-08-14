package main

import (
	"time"

	"github.com/blend/go-sdk/logger"
)

func main() {
	all := logger.MustNew(logger.OptAll())
	go func(log logger.Log) {
		ticker := time.Tick(500 * time.Millisecond)
		for {
			<-ticker
			log.Infof("this is foo")
		}
	}(all.WithPath("foo"))

	go func(log logger.Log) {
		ticker := time.Tick(500 * time.Millisecond)
		for {
			<-ticker
			log.Infof("this is bar")
		}
	}(all.WithPath("bar"))

	select {}
}
