package main

import "github.com/blend/go-sdk/logger"

// F is a helper type alias.
type F = logger.Fields

func main() {
	log := logger.All()

	log.WithFields(F{"foo": "bar"}).Info("this is a test")
	log.WithFields(F{"foo": "baz", "url": "something"}).Debug("this is a test")
}
