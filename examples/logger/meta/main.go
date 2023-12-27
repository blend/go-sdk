/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import "github.com/blend/go-sdk/logger"

// L is a helper type alias.
type L = logger.Labels

func main() {
	log := logger.All()

	log.WithLabels(L{"foo": "bar"}).Info("this is a test")
	log.WithLabels(L{"foo": "baz", "url": "something"}).Debug("this is a test")
}
