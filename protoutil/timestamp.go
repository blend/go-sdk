package protoutil

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

// Timestamp returns a proto timestamp.
//
// NOTE: protobuf timestamps are always transmitted as UTC.
func Timestamp(t time.Time) *timestamp.Timestamp {
	if t.IsZero() {
		return nil
	}
	return &timestamp.Timestamp{
		Seconds: int64(t.UTC().Unix()),
		Nanos:   int32(t.UTC().Nanosecond()),
	}
}

// FromTimestamp returns a time.Time.
//
// NOTE: protobuf timestamps are always transmitted as UTC.
func FromTimestamp(t *timestamp.Timestamp) time.Time {
	if t == nil {
		return time.Time{}
	}
	return time.Unix(t.Seconds, int64(t.Nanos)).UTC()
}
