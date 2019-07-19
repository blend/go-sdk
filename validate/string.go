package validate

import (
	"regexp"

	"github.com/blend/go-sdk/ex"
)

// String errors
const (
	ErrStringLengthMin ex.Class = "string should be a minimum length"
	ErrStringLengthMax ex.Class = "string should be a minimum length"
	ErrStringMatches   ex.Class = "string should match regular expression"
)

// String contains helpers for string validation.
type String string

// MinLen returns a validator that a string is a minimum length.
func (s String) MinLen(length int) Validator {
	return func() error {
		if len(s) < length {
			return Errorf(ErrStringLengthMin, "length: %d", length)
		}
		return nil
	}
}

// MaxLen returns a validator that a string is a minimum length.
func (s String) MaxLen(length int) Validator {
	return func() error {
		if len(s) > length {
			return Errorf(ErrStringLengthMax, "length: %d", length)
		}
		return nil
	}
}

// BetweenLen returns a validator that a string is a between a minimum and maximum length.
func (s String) BetweenLen(min, max int) Validator {
	return func() error {
		if len(s) < min {
			return Errorf(ErrStringLengthMin, "length: %d", min)
		}
		if len(s) > max {
			return Errorf(ErrStringLengthMax, "length: %d", max)
		}
		return nil
	}
}

// Matches returns a validator that a string matches a given regex.
func (s String) Matches(expression string) Validator {
	exp, err := regexp.Compile(expression)
	return func() error {
		if err != nil {
			return ex.New(err)
		}
		if !exp.MatchString(string(s)) {
			return Errorf(ErrStringMatches, "expression: %s", expression)
		}
		return nil
	}
}
