/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sentry

import "github.com/blend/go-sdk/logger"

// Constants
const (
	Platform	= "go"
	SDK		= "sentry.go"
	ListenerName	= "sentry"
)

var (
	// DefaultListenerFlags are the default log flags to notify Sentry for
	DefaultListenerFlags = []string{logger.Error, logger.Fatal}
)
