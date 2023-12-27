/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
