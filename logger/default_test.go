package logger

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestDefault(t *testing.T) {
	assert := assert.New(t)

	ensureLog()
	_log.Formatter = NewTextOutputFormatter(OptTextNoColor())

	sc := SubContext("foo", OptContextFields(Fields{"buff": "luff"}))
	assert.Equal([]string{"foo"}, sc.Path)
	assert.Equal("luff", sc.Fields["buff"])

	buf := new(bytes.Buffer)
	_log.Output = buf
	Infof("format %s", "test")
	assert.Contains(buf.String(), "[info] format test")

	buf = new(bytes.Buffer)
	_log.Output = buf
	Debugf("format %s", "test")
	assert.Contains(buf.String(), "[debug] format test")

	buf = new(bytes.Buffer)
	_log.Output = buf
	Warningf("format %s", "test")
	assert.Contains(buf.String(), "[warning] format test")

	buf = new(bytes.Buffer)
	_log.Output = buf
	Errorf("format %s", "test")
	assert.Contains(buf.String(), "[error] format test")

	buf = new(bytes.Buffer)
	_log.Output = buf
	Fatalf("format %s", "test")
	assert.Contains(buf.String(), "[fatal] format test")
}
