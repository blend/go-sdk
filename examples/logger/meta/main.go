/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
