/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"time"

	"github.com/blend/go-sdk/ex"
)

// String errors
const (
	ErrTimeBefore	ex.Class	= "time should be before"
	ErrTimeAfter	ex.Class	= "time should be after"
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
		if t.Value == nil {
			return Errorf(ErrTimeBefore, nil, "before: %v", before)
		}
		if t.Value.After(before) {
			return Errorf(ErrTimeBefore, *t.Value, "before: %v", before)
		}
		return nil
	}
}

// BeforeNowUTC returns a validator that a time should be before a given time.
func (t TimeValidators) BeforeNowUTC() Validator {
	return func() error {
		nowUTC := time.Now().UTC()
		if t.Value == nil {
			return Errorf(ErrTimeBefore, nil, "before: %v", nowUTC)
		}
		if t.Value.After(nowUTC) {
			return Errorf(ErrTimeBefore, *t.Value, "before: %v", nowUTC)
		}
		return nil
	}
}

// After returns a validator that a time should be after a given time.
func (t TimeValidators) After(after time.Time) Validator {
	return func() error {
		if t.Value == nil {
			return Errorf(ErrTimeAfter, nil, "after: %v", after)
		}
		if t.Value.Before(after) {
			return Errorf(ErrTimeAfter, *t.Value, "after: %v", after)
		}
		return nil
	}
}

// AfterNowUTC returns a validator that a time should be after a given time.
func (t TimeValidators) AfterNowUTC() Validator {
	return func() error {
		nowUTC := time.Now().UTC()
		if t.Value == nil {
			return Errorf(ErrTimeAfter, nil, "after: %v", nowUTC)
		}
		if t.Value.Before(nowUTC) {	// if value not after now == value is before now
			return Errorf(ErrTimeAfter, *t.Value, "after: %v", nowUTC)
		}
		return nil
	}
}

// Between returns a validator that a time should be after a given time.
func (t TimeValidators) Between(start, end time.Time) Validator {
	return func() error {
		if t.Value == nil {
			return Errorf(ErrTimeAfter, nil, "after: %v", start)
		}
		if t.Value.Before(start) {
			return Errorf(ErrTimeAfter, *t.Value, "after: %v", start)
		}
		if t.Value.After(end) {
			return Errorf(ErrTimeBefore, *t.Value, "before: %v", end)
		}
		return nil
	}
}
