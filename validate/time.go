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
func Time(value *time.Time) TimeValidators {
	return TimeValidators{value}
}

// TimeValidators implements validators for time.Time values.
type TimeValidators struct {
	Value *time.Time
}

// Before returns a validator that a time should be before a given time.
func (t TimeValidators) Before(before time.Time) Validator {
	return func() error {
		if t.Value == nil || (*t.Value).After(before) {
			return Errorf(ErrTimeBefore, "before: %v", before)
		}
		return nil
	}
}

// BeforeNowUTC returns a validator that a time should be before a given time.
func (t TimeValidators) BeforeNowUTC() Validator {
	return func() error {
		nowUTC := time.Now().UTC()
		if t.Value == nil || (*t.Value).After(nowUTC) {
			return Errorf(ErrTimeBefore, "before: %v", nowUTC)
		}
		return nil
	}
}

// After returns a validator that a time should be after a given time.
func (t TimeValidators) After(after time.Time) Validator {
	return func() error {
		if t.Value == nil || (*t.Value).Before(after) {
			return Errorf(ErrTimeAfter, "after: %v", after)
		}
		return nil
	}
}

// AfterNowUTC returns a validator that a time should be after a given time.
func (t TimeValidators) AfterNowUTC() Validator {
	return func() error {
		nowUTC := time.Now().UTC()
		if t.Value == nil || (*t.Value).Before(nowUTC) {
			return Errorf(ErrTimeAfter, "after: %v", nowUTC)
		}
		return nil
	}
}

// Between returns a validator that a time should be after a given time.
func (t TimeValidators) Between(before, after time.Time) Validator {
	return func() error {
		if t.Value == nil || (*t.Value).Before(before) {
			return Errorf(ErrTimeAfter, "after: %v", before)
		}
		if (*t.Value).After(after) {
			return Errorf(ErrTimeBefore, "before: %v", after)
		}
		return nil
	}
}
