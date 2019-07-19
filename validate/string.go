package validate

import (
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
	ErrStringLengthMin ex.Class = "string should be a minimum length"
	ErrStringLengthMax ex.Class = "string should be a minimum length"
	ErrStringMatches   ex.Class = "string should match regular expression"
	ErrStringIsUpper   ex.Class = "string should be uppercase"
	ErrStringIsLower   ex.Class = "string should be lowercase"
	ErrStringIsTitle   ex.Class = "string should be titlecase"
	ErrStringIsUUID    ex.Class = "string should be a uuid"
	ErrStringIsEmail   ex.Class = "string should be a valid email address"
	ErrStringIsURI     ex.Class = "string should be a valid uri"
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

// IsUpper returns a validator if a string is all uppercase.
func (s String) IsUpper() error {
	runes := []rune(string(s))
	for _, r := range runes {
		if !unicode.IsUpper(r) {
			return Error(ErrStringIsUpper)
		}
	}
	return nil
}

// IsLower returns a validator if a string is all lowercase.
func (s String) IsLower() error {
	runes := []rune(string(s))
	for _, r := range runes {
		if !unicode.IsLower(r) {
			return Error(ErrStringIsLower)
		}
	}
	return nil
}

// IsTitle returns a validator if a string is titlecase.
// Titlecase is defined as the output of strings.ToTitle(s).
func (s String) IsTitle() error {
	if strings.ToTitle(string(s)) == string(s) {
		return nil
	}
	return Error(ErrStringIsTitle)
}

// IsUUID returns if a string is a valid uuid.
func (s String) IsUUID() error {
	if _, err := uuid.Parse(string(s)); err != nil {
		return Error(ErrStringIsUUID)
	}
	return nil
}

// IsEmail returns if a string is a valid email address.
func (s String) IsEmail() error {
	if _, err := mail.ParseAddress(string(s)); err != nil {
		return Error(ErrStringIsEmail)
	}
	return nil
}

// IsURI returns if a string is a valid uri.
func (s String) IsURI() error {
	if _, err := url.ParseRequestURI(string(s)); err != nil {
		return Error(ErrStringIsURI)
	}
	return nil
}
