/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ansi_test

import (
	"testing"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/assert"
)

func TestColor256_Apply(t *testing.T) {
	t.Parallel()
	it := assert.New(t)

	actual := ansi.Color256Gold3Alt2.Apply("[CONFIG] Timeout:")
	expected := "\033[38;5;178m[CONFIG] Timeout:\033[0m"
	it.Equal(expected, actual)
}
