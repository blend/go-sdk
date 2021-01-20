package async

import "time"

// Latch states
const (
	LatchStopped  int32 = 0
	LatchStarting int32 = 1
	LatchResuming int32 = 2
	LatchStarted  int32 = 3
	LatchActive   int32 = 4
	LatchPausing  int32 = 5
	LatchPaused   int32 = 6
	LatchStopping int32 = 7
)

// Constants
const (
	DefaultQueueMaxWork        = 1 << 10
	DefaultInterval            = 500 * time.Millisecond
	DefaultShutdownGracePeriod = 10 * time.Second
)
