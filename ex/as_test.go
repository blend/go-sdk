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

type customError struct {
	err error
}

func (e *customError) Unwrap() error {
	return e.err
}

func (e *customError) Error() string {
	return e.err.Error()
}

func TestAs(t *testing.T) {
	// base sentinel errors to check against
	baseSentinel := &Ex{Class: Class("sentinel")}
	baseSentinelAlt := &Ex{Class: Class("sentinel alt")}

	exSentinelErr := New("sentinel")
	stdLibSentinelErr := errors.New("sentinel")
	nonPointerExErr := Ex{Class: Class("sentinel")}
	customErrWrappingEx := customError{err: exSentinelErr}
	customErrWrappingstdLibErr := customError{err: stdLibSentinelErr}
	customDeepWrappedErr := customError{
		err: fmt.Errorf(
			"Top of deep error: %w",
			fmt.Errorf(
				"Level 1 of deep error: %w",
				&customError{err: New(exSentinelErr)},
			),
		),
	}

	// Note: Multi's Unwrap method converts all errors in the error chain to *ex.Ex.
	multiErrorWithStdLibFirst := Append(stdLibSentinelErr, baseSentinelAlt)
	multiErrorWithAltFirst := Append(baseSentinelAlt, exSentinelErr)
	multiErrorWithSentinelFirst := Append(exSentinelErr, baseSentinelAlt)

	tests := []struct {
		name       string
		candidate  error
		expectedEx *Ex
	}{
		{
			name:       "stdlib error is not Ex pointer",
			candidate:  stdLibSentinelErr,
			expectedEx: nil,
		},
		{
			name:       "error from New is Ex pointer",
			candidate:  exSentinelErr,
			expectedEx: baseSentinel,
		},
		{
			name:       "non pointer Ex is resolved to Ex pointer",
			candidate:  &nonPointerExErr,
			expectedEx: baseSentinel,
		},
		{
			name:       "Handles custom errors that wrap Ex pointer",
			candidate:  &customErrWrappingEx,
			expectedEx: baseSentinel,
		},
		{
			name:       "wrapped stdlib error is no Ex pointer",
			candidate:  &customErrWrappingstdLibErr,
			expectedEx: nil,
		},
		{
			name:       "deeply nested Ex pointer is extracted",
			candidate:  &customDeepWrappedErr,
			expectedEx: baseSentinel,
		},
		{
			name:       "Returns first Ex in Multi error chain",
			candidate:  multiErrorWithAltFirst,
			expectedEx: baseSentinelAlt,
		},
		{
			name:       "Returns first Ex in Multi error chain alt",
			candidate:  multiErrorWithSentinelFirst,
			expectedEx: baseSentinel,
		},
		{
			name:       "Returns first error as Ex in Multi error chain",
			candidate:  multiErrorWithStdLibFirst,
			expectedEx: baseSentinel,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			its := assert.New(t)
			maybeEx := As(test.candidate)
			if test.expectedEx == nil {
				its.Nil(maybeEx)
				return
			}
			its.NotNil(maybeEx)
			its.Equal(test.expectedEx.Error(), maybeEx.Error())
		})
	}
}
