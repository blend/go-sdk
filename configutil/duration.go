package configutil

import (
	"context"
	"time"
)

// DurationSource is a type that can return a time.Duration value.
type DurationSource interface {
	// Duration should return a time.Duration if the source has a given value.
	// It should return nil if the value is not present.
	// It should return an error if there was a problem fetching the value.
	Duration(context.Context) (*time.Duration, error)
}

var (
	_ DurationSource = (*Duration)(nil)
)

// Duration implements value provider.
//
// If the value is zero, a nil is returned by the implementation indicating
// the value was not present.
//
// If you want 0 to be a valid value, you must use DurationPtr.
type Duration time.Duration

// Duration returns the value for a constant.
func (dc Duration) Duration(_ context.Context) (*time.Duration, error) {
	if dc > 0 {
		value := time.Duration(dc)
		return &value, nil
	}
	return nil, nil
}
