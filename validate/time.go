package validate

import (
	"time"

	"github.com/blend/go-sdk/ex"
)

// String errors
const (
	ErrTimeBefore ex.Class = "time should be before"
	ErrTimeAfter  ex.Class = "time should be after"
)

// Time validator singleton.
type Time time.Time

// Before returns a validator that a time should be before a given time.
func (t Time) Before(before time.Time) Validator {
	return func() error {
		if time.Time(t).After(before) {
			return Errorf(ErrTimeBefore, "before: %v", before)
		}
		return nil
	}
}

// BeforeNowUTC returns a validator that a time should be before a given time.
func (t Time) BeforeNowUTC() error {
	nowUTC := time.Now().UTC()
	if time.Time(t).After(nowUTC) {
		return Errorf(ErrTimeBefore, "before: %v", nowUTC)
	}
	return nil
}

// After returns a validator that a time should be after a given time.
func (t Time) After(after time.Time) Validator {
	return func() error {
		if time.Time(t).Before(after) {
			return Errorf(ErrTimeAfter, "after: %v", after)
		}
		return nil
	}
}

// AfterNowUTC returns a validator that a time should be after a given time.
func (t Time) AfterNowUTC() error {
	nowUTC := time.Now().UTC()
	if time.Time(t).Before(nowUTC) {
		return Errorf(ErrTimeAfter, "after: %v", nowUTC)
	}
	return nil
}

// Between returns a validator that a time should be after a given time.
func (t Time) Between(before, after time.Time) Validator {
	return func() error {
		if time.Time(t).Before(before) {
			return Errorf(ErrTimeAfter, "after: %v", before)
		}
		if time.Time(t).After(after) {
			return Errorf(ErrTimeBefore, "before: %v", after)
		}
		return nil
	}
}
