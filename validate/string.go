package validate

import (
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/blend/go-sdk/uuid"

	"github.com/blend/go-sdk/ex"
)

// String errors
const (
	ErrStringLength    ex.Class = "string should be a given length"
	ErrStringLengthMin ex.Class = "string should be a minimum length"
	ErrStringLengthMax ex.Class = "string should be a minimum length"
	ErrStringMatches   ex.Class = "string should match regular expression"
	ErrStringIsUpper   ex.Class = "string should be uppercase"
	ErrStringIsLower   ex.Class = "string should be lowercase"
	ErrStringIsTitle   ex.Class = "string should be titlecase"
	ErrStringIsUUID    ex.Class = "string should be a uuid"
	ErrStringIsEmail   ex.Class = "string should be a valid email address"
	ErrStringIsURI     ex.Class = "string should be a valid uri"
	ErrStringIsIP      ex.Class = "string should be a valid ip address"
)

// String contains helpers for string validation.
func String(value *string) StringValidators {
	return StringValidators{value}
}

// StringValidators returns string validators.
type StringValidators struct {
	Value *string
}

// MinLen returns a validator that a string is a minimum length.
// If the string is unset (nil) it will fail.
func (s StringValidators) MinLen(length int) Validator {
	return func() error {
		if s.Value == nil || len(*s.Value) < length { //if it's unset, it should fail the minimum check.
			return Errorf(ErrStringLengthMin, "length: %d", length)
		}
		return nil
	}
}

// MaxLen returns a validator that a string is a minimum length.
func (s StringValidators) MaxLen(length int) Validator {
	return func() error {
		if s.Value == nil {
			return nil
		}
		if len(*s.Value) > length {
			return Errorf(ErrStringLengthMax, "length: %d", length)
		}
		return nil
	}
}

// Length returns a validator that a string is a minimum length.
// It will fail if the value is unset (nil).
func (s StringValidators) Length(length int) Validator {
	return func() error {
		if s.Value == nil || len(*s.Value) != length {
			return Errorf(ErrStringLength, "length: %d", length)
		}
		return nil
	}
}

// BetweenLen returns a validator that a string is a between a minimum and maximum length.
// If the string is unset (nil) it will fail the minimum check.
func (s StringValidators) BetweenLen(min, max int) Validator {
	return func() error {
		if s.Value == nil || len(*s.Value) < min {
			return Errorf(ErrStringLengthMin, "length: %d", min)
		}
		if len(*s.Value) > max {
			return Errorf(ErrStringLengthMax, "length: %d", max)
		}
		return nil
	}
}

// Matches returns a validator that a string matches a given regex.
// If the value is unset it will match, i.e. not fail validation.
func (s StringValidators) Matches(expression string) Validator {
	exp, err := regexp.Compile(expression)
	return func() error {
		if err != nil {
			return ex.New(err)
		}
		if s.Value == nil {
			return nil
		}
		if !exp.MatchString(string(*s.Value)) {
			return Errorf(ErrStringMatches, "expression: %s", expression)
		}
		return nil
	}
}

// IsUpper returns a validator if a string is all uppercase.
func (s StringValidators) IsUpper() error {
	if s.Value == nil {
		return nil
	}
	runes := []rune(string(*s.Value))
	for _, r := range runes {
		if !unicode.IsUpper(r) {
			return Error(ErrStringIsUpper)
		}
	}
	return nil
}

// IsLower returns a validator if a string is all lowercase.
func (s StringValidators) IsLower() error {
	if s.Value == nil {
		return nil
	}
	runes := []rune(string(*s.Value))
	for _, r := range runes {
		if !unicode.IsLower(r) {
			return Error(ErrStringIsLower)
		}
	}
	return nil
}

// IsTitle returns a validator if a string is titlecase.
// Titlecase is defined as the output of strings.ToTitle(s).
func (s StringValidators) IsTitle() error {
	if s.Value == nil {
		return nil
	}
	if strings.ToTitle(string(*s.Value)) == string(*s.Value) {
		return nil
	}
	return Error(ErrStringIsTitle)
}

// IsUUID returns if a string is a valid uuid.
func (s StringValidators) IsUUID() error {
	if s.Value == nil {
		return nil
	}
	if _, err := uuid.Parse(string(*s.Value)); err != nil {
		return Error(ErrStringIsUUID)
	}
	return nil
}

// IsEmail returns if a string is a valid email address.
func (s StringValidators) IsEmail() error {
	if s.Value == nil {
		return nil
	}
	if _, err := mail.ParseAddress(string(*s.Value)); err != nil {
		return Error(ErrStringIsEmail)
	}
	return nil
}

// IsURI returns if a string is a valid uri.
func (s StringValidators) IsURI() error {
	if s.Value == nil {
		return nil
	}
	if _, err := url.ParseRequestURI(string(*s.Value)); err != nil {
		return Error(ErrStringIsURI)
	}
	return nil
}

// IsIP returns if a string is a valid ip address.
func (s StringValidators) IsIP() error {
	if s.Value == nil {
		return nil
	}
	if addr := net.ParseIP(string(*s.Value)); addr == nil {
		return Error(ErrStringIsIP)
	}
	return nil
}
