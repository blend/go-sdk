/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestShimLogger(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	log, err := New(
		OptOutput(buf),
		OptAll(),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	assert.Nil(err)
	defer log.Close()

	sw := NewShimWriter(log,
		OptShimWriterEventProvider(
			ShimWriterMessageEventProvider("shim"),
		),
	)
	fmt.Fprintln(sw, "this is a test")
	fmt.Fprintln(sw, "this is also a test")

	assert.NotEmpty(buf.String())
	assert.Equal("[shim] this is a test\n[shim] this is also a test\n", buf.String())
}
