/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package migration

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestAppliedEvent(t *testing.T) {
	a := assert.New(t)
	e := NewEvent(StatApplied, "Create table table_test", "label_1")
	var b bytes.Buffer
	e.WriteText(logger.NewTextOutputFormatter(), &b)
	a.Equal("\x1b[0;90m--\x1b[0m \x1b[0;34mapplied\x1b[0m label_1 \x1b[0;90m--\x1b[0m Create table table_test", b.String())
	jBytes, err := json.Marshal(e.Decompose())
	a.Nil(err)
	json := string(jBytes)
	a.Contains(json, `"result":"applied"`)
	a.Contains(json, `"body":"Create table table_test"`)
}

func TestSkippedEvent(t *testing.T) {
	a := assert.New(t)
	e := NewEvent(StatSkipped, "Create table table_test", "label_1")
	var b bytes.Buffer
	e.WriteText(logger.NewTextOutputFormatter(), &b)
	a.Equal("\x1b[0;90m--\x1b[0m \x1b[0;33mskipped\x1b[0m label_1 \x1b[0;90m--\x1b[0m Create table table_test", b.String())
	jBytes, err := json.Marshal(e.Decompose())
	a.Nil(err)
	json := string(jBytes)
	a.Contains(json, `"result":"skipped"`)
	a.Contains(json, `"body":"Create table table_test"`)
}

func TestFailedEvent(t *testing.T) {
	a := assert.New(t)
	e := NewEvent(StatFailed, "Create table table_test", "label_1")
	var b bytes.Buffer
	e.WriteText(logger.NewTextOutputFormatter(), &b)
	a.Equal("\x1b[0;90m--\x1b[0m \x1b[0;31mfailed\x1b[0m label_1 \x1b[0;90m--\x1b[0m Create table table_test", b.String())
	jBytes, err := json.Marshal(e.Decompose())
	a.Nil(err)
	json := string(jBytes)
	a.Contains(json, `"result":"failed"`)
	a.Contains(json, `"body":"Create table table_test"`)
}

func TestStatsEvent(t *testing.T) {
	a := assert.New(t)
	se := NewStatsEvent(5, 2, 0, 7)
	var b bytes.Buffer
	se.WriteText(logger.NewTextOutputFormatter(), &b)
	a.Equal("\x1b[0;32m5\x1b[0m applied \x1b[0;92m2\x1b[0m skipped \x1b[0;31m0\x1b[0m failed \x1b[0;97m7\x1b[0m total", b.String())
	jBytes, err := json.Marshal(se.Decompose())
	a.Nil(err)
	json := string(jBytes)
	a.Contains(json, `"applied":5`)
	a.Contains(json, `"skipped":2`)
	a.Contains(json, `"failed":0`)
	a.Contains(json, `"total":7`)
}
