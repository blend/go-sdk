package logger

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestMessageEvent(t *testing.T) {
	assert := assert.New(t)

	me := NewMessageEvent("flag", "an-message",
		OptMessageText("event-message"),
		OptMessageElapsed(time.Second),
	)
	assert.Equal("flag", me.Flag)
	assert.Equal("event-message", me.Text)
	assert.Equal(time.Second, me.Elapsed)

	buf := new(bytes.Buffer)
	noColor := TextOutputFormatter{
		NoColor: true,
	}

	me.WriteText(noColor, buf)
	assert.Equal("event-message (1s)", buf.String())

	contents, err := json.Marshal(me)
	assert.Nil(err)
	assert.Contains(string(contents), "event-message")
}
