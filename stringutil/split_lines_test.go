package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSplitLines(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Input    string
		Expected []string
	}{
		{"", nil},
		{"this", []string{"this"}},
		{"this\nthat", []string{"this", "that"}},
		{"this\nthat\n", []string{"this", "that"}},
		{"this\nthat\nthose", []string{"this", "that", "those"}},
		{"this\nthat\nthose\n", []string{"this", "that", "those"}},
		{"this\rthat\nthose\n", []string{"this", "that", "those"}},
		{"this\rthat\rthose\n", []string{"this", "that", "those"}},
		{"this\rthat\rthose\r", []string{"this", "that", "those"}},
		{"this\r\nthat\rthose\r", []string{"this\r", "that\rthose\r"}},
		{"this\r\nthat\r\nthose\r", []string{"this\r", "that\r", "those\r"}},
		{"this\r\nthat\r\nthose\r\n", []string{"this\r", "that\r", "those\r"}},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, SplitLines(tc.Input))
	}
}
