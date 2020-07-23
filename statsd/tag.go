package statsd

import "github.com/blend/go-sdk/stats"

// Tag is an alias / wrapper to stats.Tag
func Tag(k, v string) string {
	return stats.Tag(k, v)
}
