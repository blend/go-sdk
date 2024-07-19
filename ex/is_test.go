/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
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
	multi := Multi{New("not test class"), Class("test class")}
	testCases := []struct {
		Name     string
		Err      error
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
			Name:     "True for comparing a ClassProvider against sentinel error",
			Err:      classProvider{ErrClass: stdLibErr},
			Cause:    stdLibErr,
			Expected: true,
		},
		{
			Name:     "True for comparing a ClassProvider that implements Unwrap() against sentinel error",
			Err:      wrappingClassProvider{error: stdLibErr, ErrClass: New("another class")},
			Cause:    stdLibErr,
			Expected: true,
		},
		{
			Name:     "False for comparing multi ex against multi (due to Go slice comparison limitation)",
			Err:      Multi([]error{New("not test class"), Class("test class")}),
			Cause:    Multi([]error{New("not test class"), Class("test class")}),
			Expected: false,
		},
		{
			Name:     "False for identity comparison of multi ex (due to Go slice comparison limitation)",
			Err:      multi,
			Cause:    multi,
			Expected: false,
		},
		{
			Name: "True for comparing multi ex against class",
			Err: func() error {
				var err error = Multi{Class("test class")}
				for i := 0; i < 50; i++ {
					err = Multi{New(err), New("not test class")}
				}
				return err
			}(),
			Cause:    Class("test class"),
			Expected: true,
		},
		{
			Name: "True for comparing multi ex (created via ex.Append()) against class",
			Err: func() error {
				var err error
				for i := 0; i < 50; i++ {
					err = Append(err, New("not test class"))
				}
				return Append(err, Class("test class"))
			}(),
			Cause:    Class("test class"),
			Expected: true,
		},
		{
			Name:     "False for comparing a ClassProvider without Unwrap() against sentinel error",
			Err:      classProvider{error: stdLibErr, ErrClass: New("another class")},
			Cause:    stdLibErr,
			Expected: false,
		},
		{
			Name:     "False for comparing multi ex against different multi #1",
			Err:      Multi([]error{New("not test class"), Class("test class"), Class("?? class")}),
			Cause:    Multi([]error{New("not test class"), Class("test class")}),
			Expected: false,
		},
		{
			Name:     "False for comparing multi ex against different multi #2",
			Err:      Multi([]error{New("not test class"), Class("test class")}),
			Cause:    Multi([]error{New("test class"), Class("not test class")}),
			Expected: false,
		},
		{
			Name: "False for comparing multi ex (created via ex.Append()) against class",
			Err: func() error {
				var err error
				for i := 0; i < 50; i++ {
					err = Append(err, New("not test class"))
				}
				return Append(err, Class("test class"))
			}(),
			Cause:    Class("neither"),
			Expected: false,
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
			Name:     "True for identity comparison of native error",
			Err:      stdLibErr,
			Cause:    stdLibErr,
			Expected: true,
		},
		{
			Name:     "False for comparison between two different instances of the same native error",
			Err:      stdLibErr,
			Cause:    errors.New("sentinel"),
			Expected: false,
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

type wrappingClassProvider classProvider // classProvider is declared in util_test.go

// Unwrap implements the error interface.
func (wcp wrappingClassProvider) Unwrap() error {
	return wcp.error
}
