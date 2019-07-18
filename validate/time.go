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

// TimeValidator is a validator for strings.
type TimeValidator func(time.Time) error

// Time validator singleton.
var Time timeValidators

// timeValidators contains helpers for string validation.
type timeValidators struct{}

// Before returns a validator that a time should be before a given time.
func (t timeValidators) Before(before time.Time) TimeValidator {
	return func(t0 time.Time) error {
		if t0.After(before) {
			return Errorf(ErrTimeBefore, "before: %v", before)
		}
		return nil
	}
}

// Before returns a validator that a time should be before a given time.
func (t timeValidators) BeforeNowUTC() TimeValidator {
	return func(t0 time.Time) error {
		nowUTC := time.Now().UTC()
		if t0.After(nowUTC) {
			return Errorf(ErrTimeBefore, "before: %v", nowUTC)
		}
		return nil
	}
}

// After returns a validator that a time should be after a given time.
func (t timeValidators) After(after time.Time) TimeValidator {
	return func(t0 time.Time) error {
		if t0.Before(after) {
			return Errorf(ErrTimeAfter, "after: %v", after)
		}
		return nil
	}
}

// After returns a validator that a time should be after a given time.
func (t timeValidators) AfterNowUTC() TimeValidator {
	return func(t0 time.Time) error {
		nowUTC := time.Now().UTC()
		if t0.Before(nowUTC) {
			return Errorf(ErrTimeAfter, "after: %v", nowUTC)
		}
		return nil
	}
}

// After returns a validator that a time should be after a given time.
func (t timeValidators) Between(before, after time.Time) TimeValidator {
	return func(t0 time.Time) error {
		if t0.After(before) {
			return Errorf(ErrTimeBefore, "before: %v", before)
		}
		if t0.Before(after) {
			return Errorf(ErrTimeAfter, "after: %v", before)
		}
		return nil
	}
}
