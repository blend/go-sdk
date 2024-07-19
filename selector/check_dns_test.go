/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_CheckDNS(t *testing.T) {
	its := assert.New(t)

	invalidInputs := []string{
		"",
		"FOO",
		"invalid!",
		"!invalid",
		"inval!d",
		"-prefix",
		"suffix-",
		".dots",
		"dots.",
		"dots..dots",
		"dots-.dots",
		"dots.-dots",
		"dots-.-dots",
		"dots-.-dots",
	}
	for _, input := range invalidInputs {
		its.NotNil(CheckDNS(input), "input:", input)
	}

	validInputs := []string{
		"foo",
		"foo.bar",
		"foo-bar.moo",
		"foo-bar.moo-bar",
		"foo-bar.moo-bar",
	}
	for _, input := range validInputs {
		its.Nil(CheckDNS(input), "input:", input)
	}
}
