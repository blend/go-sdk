package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/uuid"
)

func TestParseBool(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Input    string
		Expected bool
		Err      error
	}{
		{"true", true, nil},
		{"t", true, nil},
		{"yes", true, nil},
		{"y", true, nil},
		{"1", true, nil},
		{"enabled", true, nil},
		{"on", true, nil},

		{"false", true, nil},
		{"f", true, nil},
		{"no", true, nil},
		{"n", true, nil},
		{"0", true, nil},
		{"disabled", true, nil},
		{"off", true, nil},

		{"foo", false, ErrInvalidBoolValue},
		{"", false, ErrInvalidBoolValue},
		{"00", false, ErrInvalidBoolValue},
		{uuid.V4().String(), false, ErrInvalidBoolValue},
	}

	var boolValue bool
	var err error
	for _, tc := range testCases {
		boolValue, err = ParseBool(tc.Input)
		if tc.Err != nil {
			assert.Equal(tc.Err, ex.ErrClass(err))
		} else {
			assert.Equal(tc.Expected, boolValue)
		}
	}
}
