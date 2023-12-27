/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"time"

	"github.com/blend/go-sdk/logger"
)

func main() {
	log := logger.MustNew(logger.OptConfigFromEnv())
	tick := time.Tick(time.Second)
	for range tick {
		log.Infof("it's %s", time.Now().Format(time.RFC3339))
	}
}
