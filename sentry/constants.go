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
