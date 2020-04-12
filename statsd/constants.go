package statsd

import "time"

// Constants
const (
	DefaultDialTimeout   = time.Second
	DefaultMaxPacketSize = 1 << 12 // 2^12 or 4kB
	DefaultMaxBufferSize = 32
)
