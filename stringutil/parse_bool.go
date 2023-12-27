/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

import (
	"strings"

	"github.com/blend/go-sdk/ex"
)

// Error Constants
const (
	ErrInvalidBoolValue ex.Class = "invalid bool value"
)

// MustParseBool parses a boolean value and panics if there is an error.
func MustParseBool(str string) bool {
	boolValue, err := ParseBool(str)
	if err != nil {
		panic(err)
	}
	return boolValue
}

// ParseBool parses a given string as a boolean value.
func ParseBool(str string) (bool, error) {
	strLower := strings.ToLower(strings.TrimSpace(str))
	switch strLower {
	case "true", "t", "1", "yes", "y", "enabled", "on":
		return true, nil
	case "false", "f", "0", "no", "n", "disabled", "off":
		return false, nil
	}
	return false, ex.New(ErrInvalidBoolValue, ex.OptMessage(str))
}
