/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"errors"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestIs(t *testing.T) {
	stdLibErr := errors.New("sentinel")
	testCases := []struct {
		Name     string
		Err      interface{}
		Cause    error
		Expected bool
	}{
		{
			Name:     "True for equal classes",
			Err:      Class("test class"),
			Cause:    Class("test class"),
			Expected: true,
		},
		{
			Name:     "True for comparing ex against class",
			Err:      New("test class"),
			Cause:    Class("test class"),
			Expected: true,
		},
		{
			Name:     "True for comparing ex against ex with equal class",
			Err:      New("test class"),
			Cause:    New("test class"),
			Expected: true,
		},
		{
			Name:     "True for comparing ex against ex with equal sentinel errors",
			Err:      New(stdLibErr),
			Cause:    stdLibErr,
			Expected: true,
		},
		{
			Name:     "True for comparing ex with wrapped sentinel error against sentinel error",
			Err:      New(fmt.Errorf("outer err: %w", stdLibErr)),
			Cause:    stdLibErr,
			Expected: true,
		},
		{
			Name:     "True for comparing multi ex against class",
			Err:      Multi([]error{New("test class"), Class("not test class")}),
			Cause:    Class("not test class"),
			Expected: true,
		},
		{
			Name:     "False for different classes",
			Err:      Class("not test class"),
			Cause:    New("test class"),
			Expected: false,
		},
		{
			Name:     "False for ex with different classes",
			Err:      New("test class"),
			Cause:    New("not test class"),
			Expected: false,
		},
		{
			Name:     "False for ex comparison against nil",
			Err:      New("test class"),
			Cause:    nil,
			Expected: false,
		},
		{
			Name:     "True for wrapped sentinel error",
			Err:      fmt.Errorf("outer err: %w", stdLibErr),
			Cause:    stdLibErr,
			Expected: true,
		},
		{
			Name:     "False for nil comparisons",
			Err:      nil,
			Cause:    nil,
			Expected: false,
		},
		{
			Name:     "False for nil error against class",
			Err:      nil,
			Cause:    Class("test class"),
			Expected: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			its := assert.New(t)
			its.Equal(tc.Expected, Is(tc.Err, tc.Cause))
		})
	}
}
