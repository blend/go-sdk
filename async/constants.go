/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import (
	"time"

	"github.com/blend/go-sdk/ex"
)

// Latch states
const (
	LatchStopped	int32	= 0
	LatchStarting	int32	= 1
	LatchResuming	int32	= 2
	LatchStarted	int32	= 3
	LatchActive	int32	= 4
	LatchPausing	int32	= 5
	LatchPaused	int32	= 6
	LatchStopping	int32	= 7
)

// Constants
const (
	DefaultQueueMaxWork		= 1 << 10
	DefaultInterval			= 500 * time.Millisecond
	DefaultShutdownGracePeriod	= 10 * time.Second
)

// Errors
var (
	ErrCannotStart	ex.Class	= "cannot start; already started"
	ErrCannotStop	ex.Class	= "cannot stop; already stopped"
	ErrCannotCancel	ex.Class	= "cannot cancel; already canceled"
)
