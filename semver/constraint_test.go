/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package semver

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewConstraint(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		input string
		count int
		err   bool
	}{
		{">= 1.2", 1, false},
		{"1.0", 1, false},
		{">= 1.x", 0, true},
		{">= 1.2, < 1.0", 2, false},

		// Out of bounds
		{"11387778780781445675529500000000000000000", 0, true},
	}

	for _, tc := range cases {
		v, err := NewConstraint(tc.input)
		assert.False(!tc.err && err != nil, fmt.Sprintf("error for input %s: %s", tc.input, err))
		assert.False(tc.err && err == nil, fmt.Sprintf("expected error for input: %s", tc.input))
		assert.Len(v, tc.count, fmt.Sprintf("input: %s\nexpected len: %d\nactual: %d", tc.input, tc.count, len(v)))
	}
}

func TestConstraintCheck(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		constraint string
		version    string
		check      bool
	}{
		{">= 1.0, < 1.2", "1.1.5", true},
		{"< 1.0, < 1.2", "1.1.5", false},
		{"= 1.0", "1.1.5", false},
		{"= 1.0", "1.0.0", true},
		{"1.0", "1.0.0", true},
		{"~> 1.0", "2.0", false},
		{"~> 1.0", "1.1", true},
		{"~> 1.0", "1.2.3", true},
		{"~> 1.0.0", "1.2.3", false},
		{"~> 1.0.0", "1.0.7", true},
		{"~> 1.0.0", "1.1.0", false},
		{"~> 1.0.7", "1.0.4", false},
		{"~> 1.0.7", "1.0.7", true},
		{"~> 1.0.7", "1.0.8", true},
		{"~> 1.0.7", "1.0.7.5", true},
		{"~> 1.0.7", "1.0.6.99", false},
		{"~> 1.0.7", "1.0.8.0", true},
		{"~> 1.0.9.5", "1.0.9.5", true},
		{"~> 1.0.9.5", "1.0.9.4", false},
		{"~> 1.0.9.5", "1.0.9.6", true},
		{"~> 1.0.9.5", "1.0.9.5.0", true},
		{"~> 1.0.9.5", "1.0.9.5.1", true},
		{"~> 2.0", "2.1.0-beta", false},
		{"~> 2.1.0-a", "2.2.0", false},
		{"~> 2.1.0-a", "2.1.0", false},
		{"~> 2.1.0-a", "2.1.0-beta", true},
		{"~> 2.1.0-a", "2.2.0-alpha", false},
		{"> 2.0", "2.1.0-beta", false},
		{">= 2.1.0-a", "2.1.0-beta", true},
		{">= 2.1.0-a", "2.1.1-beta", false},
		{">= 2.0.0", "2.1.0-beta", false},
		{">= 2.1.0-a", "2.1.1", true},
		{">= 2.1.0-a", "2.1.1-beta", false},
		{">= 2.1.0-a", "2.1.0", true},
		{"<= 2.1.0-a", "2.0.0", true},
	}

	for _, tc := range cases {
		c, err := NewConstraint(tc.constraint)
		assert.Nil(err)

		v, err := NewVersion(tc.version)
		assert.Nil(err)

		actual := c.Check(v)
		expected := tc.check
		assert.Equal(expected, actual)
	}
}

func TestConstraintsString(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		constraint string
		result     string
	}{
		{">= 1.0, < 1.2", ""},
		{"~> 1.0.7", ""},
	}

	for _, tc := range cases {
		c, err := NewConstraint(tc.constraint)
		assert.Nil(err)

		actual := c.String()
		expected := tc.result
		if expected == "" {
			expected = tc.constraint
		}

		assert.Equal(expected, actual)
	}
}
