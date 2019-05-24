package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMessageEvent(t *testing.T) {
	assert := assert.New(t)

	me := NewMessageEvent("flag", "event-message")
	assert.Equal("flag", me.Flag)
	assert.Equal("event-message", me.Message)

	buf := new(bytes.Buffer)
	noColor := TextOutputFormatter{
		NoColor: true,
	}

	me.WriteText(noColor, buf)
	assert.Equal("event-message", buf.String())

	contents, err := json.Marshal(me)
	assert.Nil(err)
	assert.Contains(string(contents), "event-message")
}
