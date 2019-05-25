package logger

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestTimedEvent(t *testing.T) {
	assert := assert.New(t)

	tme := Timedf("flag", time.Second, "event-message")
	assert.Equal("flag", tme.Flag)
	assert.Equal(time.Second, tme.Elapsed)
	assert.Equal("event-message", tme.Message)

	buf := new(bytes.Buffer)
	noColor := TextOutputFormatter{
		NoColor: true,
	}

	tme.WriteText(noColor, buf)
	assert.Equal("event-message (1s)", buf.String())

	contents, err := json.Marshal(tme)
	assert.Nil(err)
	assert.Contains(string(contents), "event-message")
}
