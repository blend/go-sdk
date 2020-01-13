package logger

import (
	"bytes"
	"context"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestJSONOutputFormatter(t *testing.T) {
	assert := assert.New(t)

	jf := NewJSONOutputFormatter(
		OptJSONPretty(),
		OptJSONPrettyPrefix("    "),
		OptJSONPrettyIndent("\t\t"),
	)
	assert.True(jf.Pretty)
	assert.Equal("    ", jf.PrettyPrefixOrDefault())
	assert.Equal("\t\t", jf.PrettyIndentOrDefault())
	jf.Pretty = false

	me := NewMessageEvent(Info, "this is a test")

	buf := new(bytes.Buffer)
	assert.Nil(jf.WriteFormat(context.Background(), buf, me))

	assert.Contains(buf.String(), "\"text\":\"this is a test\"")

	jf.Pretty = true
	jf.PrettyPrefix = ""
	jf.PrettyIndent = "\t"

	buf = new(bytes.Buffer)
	assert.Nil(jf.WriteFormat(context.Background(), buf, me))
	assert.Contains(buf.String(), "\t\"text\": \"this is a test\"\n")
}
