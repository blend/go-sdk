package cron

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOutputListeners(t *testing.T) {
	assert := assert.New(t)

	lines := make(map[string][]OutputLine)

	lw := new(OutputListeners)
	lw.Add("one", func(l OutputLine) {
		lines["one"] = append(lines["one"], l)
	})
	lw.Add("two", func(l OutputLine) {
		lines["two"] = append(lines["two"], l)
	})

	lw.Trigger(OutputLine{Data: []byte("test1")})
	lw.Trigger(OutputLine{Data: []byte("test2")})
	assert.Len(lines["one"], 2)
	assert.Len(lines["two"], 2)

	lw.Remove("two")

	lw.Trigger(OutputLine{Data: []byte("test2")})
	assert.Len(lines["one"], 3)
	assert.Len(lines["two"], 2)

	removeTest := new(OutputListeners)
	removeTest.Remove("something")
}
