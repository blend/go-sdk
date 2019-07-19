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

// StringValidator is a validator for strings.
type StringValidator func(string) error

// String validator singleton.
var String stringValidators

// string contains helpers for string validation.
type stringValidators struct{}

// Min returns a validator that a string is a minimum length.
func (s stringValidators) Min(length int) StringValidator {
	return func(v string) error {
		if len(v) < length {
			return Errorf(ErrStringLengthMin, "length: %d", length)
		}
		return nil
	}
}

// Max returns a validator that a string is a minimum length.
func (s stringValidators) Max(length int) StringValidator {
	return func(v string) error {
		if len(v) > length {
			return Errorf(ErrStringLengthMax, "length: %d", length)
		}
		return nil
	}
}

// Between returns a validator that a string is a between a minimum and maximum length.
func (s stringValidators) Between(min, max int) StringValidator {
	return func(v string) error {
		if len(v) < min {
			return Errorf(ErrStringLengthMin, "length: %d", min)
		}
		if len(v) > max {
			return Errorf(ErrStringLengthMax, "length: %d", max)
		}
		return nil
	}
}

// Min returns a validator that a string is a minimum length.
func (s stringValidators) Matches(expression string) StringValidator {
	exp, err := regexp.Compile(expression)
	return func(v string) error {
		if err != nil {
			return ex.New(err)
		}
		if !exp.MatchString(v) {
			return Errorf(ErrStringMatches, "expression: %s", expression)
		}
		return nil
	}
}
