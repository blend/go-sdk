package statsd

import "strings"

// Tag formats a tag with a given key and value.
func Tag(key, value string) string {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(key)
	return key + ":" + value
}
