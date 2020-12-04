package ex

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestIs(t *testing.T) {
	it := assert.New(t)

	testCases := []struct {
		Err      interface{}
		Cause    error
		Expected bool
	}{
		{Err: Class("test class"), Cause: Class("test class"), Expected: true},
		{Err: New("test class"), Cause: Class("test class"), Expected: true},
		{Err: New("test class"), Cause: New("test class"), Expected: true},
		{Err: Multi([]error{New("test class"), Class("not test class")}), Cause: Class("not test class"), Expected: true},
		{Err: Class("not test class"), Cause: New("test class"), Expected: false},
		{Err: New("test class"), Cause: New("not test class"), Expected: false},
		{Err: New("test class"), Cause: nil, Expected: false},
		{Err: nil, Cause: nil, Expected: false},
		{Err: nil, Cause: Class("test class"), Expected: false},
	}

	for index, tc := range testCases {
		it.Equal(tc.Expected, Is(tc.Err, tc.Cause), fmt.Sprintf("test case %d", index))
	}
}
