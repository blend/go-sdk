/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestStdlibShim(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	log, err := New(
		OptOutput(buf),
		OptAll(),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	assert.Nil(err)
	defer log.Close()

	shim := StdlibShim(log, OptShimWriterEventProvider(ShimWriterErrorEventProvider("error")))

	shim.Println("this is a test")
	shim.Println("this is another test")

	assert.NotEmpty(buf.String())
	assert.Equal("[error] this is a test\n[error] this is another test\n", buf.String())
}
