package logger

import (
	"bytes"
	"context"
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
	defer log.Close()
	assert.Nil(err)

	shim := StdlibShim(context.Background(), "errors", log)

	shim.Println("this is a test")
	shim.Println("this is another test")

	assert.NotEmpty(buf.String())
	assert.Equal("[errors] this is a test\n[errors] this is another test\n", buf.String())
}
