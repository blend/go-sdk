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
	ErrStringRequired  ex.Class = "string should be set"
	ErrStringLength    ex.Class = "string should be a given length"
	ErrStringLengthMin ex.Class = "string should be a minimum length"
	ErrStringLengthMax ex.Class = "string should be a maximum length"
	ErrStringMatches   ex.Class = "string should match regular expression"
	ErrStringIsUpper   ex.Class = "string should be uppercase"
	ErrStringIsLower   ex.Class = "string should be lowercase"
	ErrStringIsTitle   ex.Class = "string should be titlecase"
	ErrStringIsUUID    ex.Class = "string should be a uuid"
	ErrStringIsEmail   ex.Class = "string should be a valid email address"
	ErrStringIsURI     ex.Class = "string should be a valid uri"
	ErrStringIsIP      ex.Class = "string should be a valid ip address"
	ErrStringIsSlug    ex.Class = "string should be a valid slug (i.e. matching [0-9,a-z,A-Z,_,-])"
	ErrStringIsOneOf   ex.Class = "string should be one of a set of values"
)

// String contains helpers for string validation.
func String(value *string) StringValidators {
	return StringValidators{value}
}

// StringValidators returns string validators.
type StringValidators struct {
	Value *string
}

// Required returns a validator that a string is set and not zero length.
func (s StringValidators) Required() Validator {
	return func() error {
		if s.Value == nil {
			return Error(ErrStringRequired, nil)
		}
		if len(*s.Value) == 0 {
			return Error(ErrStringRequired, nil)
		}
		return nil
	}
}

// MinLen returns a validator that a string is a minimum length.
// If the string is unset (nil) it will fail.
func (s StringValidators) MinLen(length int) Validator {
	return func() error {
		if s.Value == nil {
			return Errorf(ErrStringLengthMin, nil, "length: %d", length)
		}
		if len(*s.Value) < length { //if it's unset, it should fail the minimum check.
			return Errorf(ErrStringLengthMin, *s.Value, "length: %d", length)
		}
		return nil
	}
}

// MaxLen returns a validator that a string is a minimum length.
// It will pass if the string is unset (nil).
func (s StringValidators) MaxLen(length int) Validator {
	return func() error {
		if s.Value == nil {
			return nil
		}
		if len(*s.Value) > length {
			return Errorf(ErrStringLengthMax, *s.Value, "length: %d", length)
		}
		return nil
	}
}

// Length returns a validator that a string is a minimum length.
// It will error if the string is unset (nil).
func (s StringValidators) Length(length int) Validator {
	return func() error {
		if s.Value == nil {
			return Errorf(ErrStringLength, nil, "length: %d", length)
		}
		if len(*s.Value) != length {
			return Errorf(ErrStringLength, *s.Value, "length: %d", length)
		}
		return nil
	}
}

// BetweenLen returns a validator that a string is a between a minimum and maximum length.
// It will error if the string is unset (nil).
func (s StringValidators) BetweenLen(min, max int) Validator {
	return func() error {
		if s.Value == nil {
			return Errorf(ErrStringLengthMin, nil, "length: %d", min)
		}
		if len(*s.Value) < min {
			return Errorf(ErrStringLengthMin, *s.Value, "length: %d", min)
		}
		if len(*s.Value) > max {
			return Errorf(ErrStringLengthMax, *s.Value, "length: %d", max)
		}
		return nil
	}
}

// Matches returns a validator that a string matches a given regex.
// It will error if the string is unset (nil).
func (s StringValidators) Matches(expression string) Validator {
	exp, err := regexp.Compile(expression)
	return func() error {
		if err != nil {
			return ex.New(err)
		}
		if s.Value == nil {
			return Errorf(ErrStringMatches, nil, "expression: %s", expression)
		}
		if !exp.MatchString(string(*s.Value)) {
			return Errorf(ErrStringMatches, *s.Value, "expression: %s", expression)
		}
		return nil
	}
}

// IsUpper returns a validator if a string is all uppercase.
// It will error if the string is unset (nil).
func (s StringValidators) IsUpper() Validator {
	return func() error {
		if s.Value == nil {
			return Error(ErrStringIsUpper, nil)
		}
		runes := []rune(string(*s.Value))
		for _, r := range runes {
			if !unicode.IsUpper(r) {
				return Error(ErrStringIsUpper, *s.Value)
			}
		}
		return nil
	}
}

// IsLower returns a validator if a string is all lowercase.
// It will error if the string is unset (nil).
func (s StringValidators) IsLower() Validator {
	return func() error {
		if s.Value == nil {
			return Error(ErrStringIsLower, nil)
		}
		runes := []rune(string(*s.Value))
		for _, r := range runes {
			if !unicode.IsLower(r) {
				return Error(ErrStringIsLower, *s.Value)
			}
		}
		return nil
	}
}

// IsTitle returns a validator if a string is titlecase.
// Titlecase is defined as the output of strings.ToTitle(s).
// It will error if the string is unset (nil).
func (s StringValidators) IsTitle() Validator {
	return func() error {
		if s.Value == nil {
			return Error(ErrStringIsTitle, nil)
		}
		if strings.ToTitle(string(*s.Value)) == string(*s.Value) {
			return nil
		}
		return Error(ErrStringIsTitle, *s.Value)
	}
}

// IsUUID returns if a string is a valid uuid.
// It will error if the string is unset (nil).
func (s StringValidators) IsUUID() Validator {
	return func() error {
		if s.Value == nil {
			return Error(ErrStringIsUUID, nil)
		}
		if _, err := uuid.Parse(string(*s.Value)); err != nil {
			return Error(ErrStringIsUUID, *s.Value)
		}
		return nil
	}
}

// IsEmail returns if a string is a valid email address.
func (s StringValidators) IsEmail() Validator {
	return func() error {
		if s.Value == nil {
			return Error(ErrStringIsEmail, nil)
		}
		if _, err := mail.ParseAddress(string(*s.Value)); err != nil {
			return Error(ErrStringIsEmail, *s.Value)
		}
		return nil
	}
}

// IsURI returns if a string is a valid uri.
// It will error if the string is unset (nil).
func (s StringValidators) IsURI() Validator {
	return func() error {
		if s.Value == nil {
			return Error(ErrStringIsURI, nil)
		}
		if _, err := url.ParseRequestURI(string(*s.Value)); err != nil {
			return Error(ErrStringIsURI, *s.Value)
		}
		return nil
	}
}

// IsIP returns if a string is a valid ip address.
// It will error if the string is unset (nil).
func (s StringValidators) IsIP() Validator {
	return func() error {
		if s.Value == nil {
			return Error(ErrStringIsIP, nil)
		}
		if addr := net.ParseIP(string(*s.Value)); addr == nil {
			return Error(ErrStringIsIP, *s.Value)
		}
		return nil
	}
}

// IsSlug returns if a string is a valid slug as defined by the match rule [0-9,a-z,A-Z,_,-].
// It will error if the string is unset (nil).
func (s StringValidators) IsSlug() Validator {
	return func() error {
		if s.Value == nil {
			return Error(ErrStringIsSlug, nil)
		}
		for _, c := range *s.Value {
			if unicode.IsLetter(c) {
				continue
			}
			if unicode.IsDigit(c) {
				continue
			}
			if c == '-' || c == '_' {
				continue
			}
			return Error(ErrStringIsSlug, *s.Value)
		}
		return nil
	}
}

// IsOneOf validates a string is one of a known set of values.
func (s StringValidators) IsOneOf(values ...string) Validator {
	return func() error {
		if s.Value == nil {
			return Error(ErrStringIsOneOf, s.Value, strings.Join(values, ", "))
		}
		for _, value := range values {
			if *s.Value == value {
				return nil
			}
		}
		return Error(ErrStringIsOneOf, s.Value, strings.Join(values, ", "))
	}
}
