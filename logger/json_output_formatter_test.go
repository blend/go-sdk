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
	assert.Empty(jf.PrettyPrefixOrDefault())
	assert.Equal("\t", jf.PrettyIndentOrDefault())
	jf.Pretty = false

	me := NewMessageEvent(Info, "this is a test")

	buf := new(bytes.Buffer)
	assert.Nil(jf.WriteFormat(context.Background(), buf, me))

	assert.Contains(buf.String(), "\"text\":\"this is a test\"")
}
