/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMultiAppend(t *testing.T) {
	it := assert.New(t)

	ex0 := New(New("hi0"))
	ex1 := New(fmt.Errorf("hi1"))
	ex2 := New("hi2")

	m := Append(ex0, ex1, ex2)

	it.True(strings.HasPrefix(m.Error(), `3 errors occurred:`), m.Error())

	it.Equal([]error{ex0, ex1, ex2}, m.(Multi))

	m = nil
	m = Append(m, ex0)
	m = Append(m, ex1)
	m = Append(m, ex2)

	it.True(strings.HasPrefix(m.Error(), `3 errors occurred:`), m.Error())

	it.Equal([]error{ex0, ex1, ex2}, m.(Multi))
}

func TestMultiUnwrap(t *testing.T) {
	it := assert.New(t)
	m := Multi{New("hi0"), New("hi1"), New("hi2")}
	it.Equal(m, m.Unwrap())
}

func TestMultiIsCompatibility(t *testing.T) {
	it := assert.New(t)
	m1 := Multi{New("hi0"), New("hi1"), New("hi2")}
	m2 := Multi{New("hi0"), New("hi1"), New("hi2")}

	it.True(errors.Is(m1, New("hi1")))
	it.False(errors.Is(m1, New("hi3")))
	it.False(errors.Is(m1, append(m2, New("hi3"))))
	it.False(errors.Is(Multi(append(m1, New("hi3"))), m2))

	// Multi to Multi comparisons always return false (Golang can't do slice comparison)
	it.False(errors.Is(m1, m1))
	it.False(errors.Is(m1, m2))
}

func TestMultiError(t *testing.T) {
	it := assert.New(t)

	err := Multi{
		New("hi0"),
		Multi{
			New("hi1.0"),
		},
		Multi{
			New("hi2.0"),
			New("hi2.1").WithMessagef("message2.1"),
			New("hi2.3").WithInner(Multi{
				New("hi2.3.0"),
				errors.New("hi2.3.1"),
			}),
			Multi{
				New("hi2.4.0"),
				errors.New("hi2.4.1"),
			},
		},
		errors.New("hi3"),
	}

	fullError := strings.Join([]string{
		"8 errors occurred:",
		"\t* hi0",
		"\t* 1 error occurred:",
		"\t\t\t* hi1.0",
		"\t* 5 errors occurred:",
		"\t\t\t* hi2.0",
		"\t\t\t* hi2.1",
		"\t\t\t* hi2.3",
		"\t\t\t* 2 errors occurred:",
		"\t\t\t\t\t* hi2.4.0",
		"\t\t\t\t\t* hi2.4.1",
		"\t* hi3",
	}, "\n")

	// Default error format
	it.Equal(fullError, err.Error())

	// With length and depth limits
	formatted, total := err.errorString(2, 1, 0)
	it.Equal(strings.Join([]string{
		"8 errors occurred:",
		"\t* hi0",
		"\t* 1 error occurred:",
		"\t\t\t... depth limit reached ...",
		"\t... and 6 more",
	}, "\n"), formatted)
	it.Equal(8, total)

	// Length and depth limits don't consider Multi errors wrapped in other kinds
	// of errors.
	wrapped := Multi{
		New(err),
	}
	formatted, total = wrapped.errorString(2, 1, 0)
	it.Equal(fmt.Sprintf("1 error occurred:\n\t* %s", indent("\t\t", fullError)), formatted)
	it.Equal(1, total)

	// Unlimited cases
	it.Equal(fullError, err.FullError())
	formatted, total = err.errorString(-1, -1, 0)
	it.Equal(fullError, formatted)
	it.Equal(8, total)
	formatted, total = err.errorString(5555, 9999, 0)
	it.Equal(fullError, formatted)
	it.Equal(8, total)

	// Edge cases for the limits
	formatted, total = err.errorString(-1, 0, 0)
	it.Equal("8 errors occurred:\n\t... depth limit reached ...", formatted)
	it.Equal(8, total)
	formatted, total = err.errorString(0, -1, 0)
	it.Equal("8 errors occurred:\n\t... and 8 more", formatted)
	it.Equal(8, total)

	// Empty Multi
	it.Equal("", Multi{}.Error())
}
