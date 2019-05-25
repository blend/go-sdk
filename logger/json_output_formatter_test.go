package logger

import (
	"bytes"
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestJSONOutputFormatter(t *testing.T) {
	assert := assert.New(t)

	jf := NewJSONOutputFormatter(OptJSONPretty())
	assert.True(jf.Pretty)
	jf.Pretty = false

	me := NewMessageEvent(Info, "this is a test")

	buf := new(bytes.Buffer)
	assert.Nil(jf.WriteFormat(context.Background(), buf, me))

	assert.Contains(buf.String(), "\"message\":\"this is a test\"")
}
