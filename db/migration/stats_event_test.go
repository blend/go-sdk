package migration

import (
	"bytes"
	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"testing"
)

func TestStatsEvent(t *testing.T) {
	a := assert.New(t)
	se := NewStatsEvent(5, 2, 0, 7)
	var b bytes.Buffer
	se.WriteText(logger.NewTextOutputFormatter(), &b)
	a.Equal("\x1b[0;32m5\x1b[0m applied \x1b[0;92m2\x1b[0m skipped \x1b[0;31m0\x1b[0m failed \x1b[0;97m7\x1b[0m total", b.String())
	jBytes, err := se.MarshalJSON()
	a.Nil(err)
	json := string(jBytes)
	a.Contains(json, `"applied":5`)
	a.Contains(json, `"skipped":2`)
	a.Contains(json, `"failed":0`)
	a.Contains(json, `"total":7`)
}
