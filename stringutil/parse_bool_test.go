/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

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
		Input		string
		Expected	bool
		Err		error
	}{
		{"true", true, nil},
		{"t", true, nil},
		{"yes", true, nil},
		{"y", true, nil},
		{"1", true, nil},
		{"enabled", true, nil},
		{"on", true, nil},

		{"false", false, nil},
		{"f", false, nil},
		{"no", false, nil},
		{"n", false, nil},
		{"0", false, nil},
		{"disabled", false, nil},
		{"off", false, nil},

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
