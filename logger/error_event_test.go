package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/assert"
)

func TestNewErrorEvent(t *testing.T) {
	assert := assert.New(t)

	/// stuff
	ee := NewErrorEvent(Fatal, fmt.Errorf("not a test"), OptErrorEventErr(fmt.Errorf("only a test")), OptErrorEventState("foo"), OptErrorEventMetaOptions(OptEventMetaFlagColor(ansi.ColorBlue)))
	assert.Equal(Fatal, ee.GetFlag())
	assert.Equal("only a test", ee.Err.Error())
	assert.Equal("foo", ee.State)
	assert.Equal(ansi.ColorBlue, ee.GetFlagColor())

	buf := new(bytes.Buffer)
	tf := TextOutputFormatter{
		NoColor: true,
	}

	ee.WriteText(tf, buf)
	assert.Equal("only a test", buf.String())

	contents, err := json.Marshal(ee)
	assert.Nil(err)
	assert.Contains(string(contents), "only a test")
}

func TestErrorEventListener(t *testing.T) {
	assert := assert.New(t)

	ee := NewErrorEvent(Fatal, fmt.Errorf("only a test"))

	var didCall bool
	ml := NewErrorEventListener(func(ctx context.Context, e *ErrorEvent) {
		didCall = true
	})

	ml(context.Background(), ee)
	assert.True(didCall)
}
