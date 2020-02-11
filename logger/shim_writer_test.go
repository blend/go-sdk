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
	defer log.Close()
	assert.Nil(err)

	sw := NewShimWriter(log,
		OptShimWriterEventProvider(
			ShimWriterMessageEventProvider("shim"),
		),
	)
	fmt.Fprintf(sw, "this is a test\n")
	fmt.Fprintf(sw, "this is also a test\n")

	assert.NotEmpty(buf.String())
	assert.Equal("[shim] this is a test\n[shim] this is also a test\n", buf.String())
}
