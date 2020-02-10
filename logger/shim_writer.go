package logger

import (
	"context"
	"strings"
)

// ShimWriter is a type that implements io.Writer with
// a logger backend.
type ShimWriter struct {
	Context context.Context
	Flag    string
	Log     Triggerable
}

// Write implements io.Writer.
func (sw ShimWriter) Write(contents []byte) (count int, err error) {
	sw.Log.Trigger(sw.Context, NewMessageEvent(sw.Flag, strings.TrimSpace(string(contents))))
	count = len(contents)
	return
}
