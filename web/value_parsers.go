package web

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Common errors
var (
	ErrInvalidBoolValue = fmt.Errorf("invalid boolean value")
)

// BoolValue parses a value as an bool.
// If the input error is set it short circuits.
func BoolValue(value string, inputErr error) (output bool, err error) {
	if inputErr != nil {
		err = inputErr
		return
	}
	switch strings.ToLower(value) {
	case "1", "true", "yes", "on":
		output = true
	case "0", "false", "no", "off":
		output = false
	default:
		err = ErrInvalidBoolValue
	}
	return
}

// IntValue parses a value as an int.
// If the input error is set it short circuits.
func IntValue(value string, inputErr error) (output int, err error) {
	if inputErr != nil {
		err = inputErr
		return
	}
	output, err = strconv.Atoi(value)
	return
}

// Int64Value parses a value as an int64.
// If the input error is set it short circuits.
func Int64Value(value string, inputErr error) (output int64, err error) {
	if inputErr != nil {
		err = inputErr
		return
	}
	output, err = strconv.ParseInt(value, 10, 64)
	return
}

// Float64Value parses a value as an float64.
// If the input error is set it short circuits.
func Float64Value(value string, inputErr error) (output float64, err error) {
	if inputErr != nil {
		err = inputErr
		return
	}
	output, err = strconv.ParseFloat(value, 64)
	return
}

// DurationValue parses a value as an time.Duration.
// If the input error is set it short circuits.
func DurationValue(value string, inputErr error) (output time.Duration, err error) {
	if inputErr != nil {
		err = inputErr
		return
	}
	output, err = time.ParseDuration(value)
	return
}

// StringValue just returns the string directly from a value error pair.
func StringValue(value string, _ error) string {
	return value
}

// CSVValue just returns the string directly from a value error pair.
func CSVValue(value string, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	return strings.Split(value, ","), nil
}
