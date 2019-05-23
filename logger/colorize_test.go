package logger

import (
	"testing"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/assert"
)

func TestFlagTextColor(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Flag     string
		Expected ansi.Color
	}{
		{Info, ansi.ColorLightWhite},
		{Debug, ansi.ColorLightYellow},
		{Warning, ansi.ColorLightYellow},
		{Error, ansi.ColorRed},
		{Fatal, ansi.ColorRed},
		{"foo", DefaultFlagTextColor},
		{"", DefaultFlagTextColor},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, FlagTextColor(tc.Flag))
	}
}
