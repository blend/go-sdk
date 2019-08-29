package stats

import (
	"github.com/blend/go-sdk/logger"
)

// MetricNames are names we use when sending data to the collectors.
const (
	MetricNameError string = string(logger.Error)
)

// Tag names are names for tags, either on metrics or traces.
const (
	TagService   string = "service"
	TagJob       string = "job"
	TagEnv       string = "env"
	TagHostname  string = "hostname"
	TagContainer string = "container"
	TagSeverity  string = "severity"
	TagError     string = "error"
	TagClass     string = "class"
)

// Specialized / default values
const (
	ListenerNameStats string = "stats"
)

// Tag creates a new tag.
func Tag(key, value string) string {
	return key + ":" + value
}
