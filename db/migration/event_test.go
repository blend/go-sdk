package migration

import (
	"bytes"
	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"testing"
)

func TestAppliedEvent(t *testing.T) {
	a := assert.New(t)
	e := NewEvent(StatApplied, "Create table table_test", "label_1")
	var b bytes.Buffer
	e.WriteText(logger.NewTextOutputFormatter(), &b)
	a.Equal("\x1b[0;90m--\x1b[0m \x1b[0;34mapplied\x1b[0m label_1 \x1b[0;90m--\x1b[0m Create table table_test", b.String())
	jBytes, err := e.MarshalJSON()
	a.Nil(err)
	json := string(jBytes)
	a.Contains(json, `"result":"applied"`)
	a.Contains(json, `"body":"Create table table_test"`)
	a.Contains(json, `"flag":"db.migration"`)
}

func TestSkippedEvent(t *testing.T) {
	a := assert.New(t)
	e := NewEvent(StatSkipped, "Create table table_test", "label_1")
	var b bytes.Buffer
	e.WriteText(logger.NewTextOutputFormatter(), &b)
	a.Equal("\x1b[0;90m--\x1b[0m \x1b[0;33mskipped\x1b[0m label_1 \x1b[0;90m--\x1b[0m Create table table_test", b.String())
	jBytes, err := e.MarshalJSON()
	a.Nil(err)
	json := string(jBytes)
	a.Contains(json, `"result":"skipped"`)
	a.Contains(json, `"body":"Create table table_test"`)
	a.Contains(json, `"flag":"db.migration"`)
}

func TestFailedEvent(t *testing.T) {
	a := assert.New(t)
	e := NewEvent(StatFailed, "Create table table_test", "label_1")
	var b bytes.Buffer
	e.WriteText(logger.NewTextOutputFormatter(), &b)
	a.Equal("\x1b[0;90m--\x1b[0m \x1b[0;31mfailed\x1b[0m label_1 \x1b[0;90m--\x1b[0m Create table table_test", b.String())
	jBytes, err := e.MarshalJSON()
	a.Nil(err)
	json := string(jBytes)
	a.Contains(json, `"result":"failed"`)
	a.Contains(json, `"body":"Create table table_test"`)
	a.Contains(json, `"flag":"db.migration"`)
}
