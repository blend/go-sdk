package bufferutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLineHandlers(t *testing.T) {
	assert := assert.New(t)

	lines := make(map[string][]Line)

	lw := new(LineHandlers)
	lw.Add("one", func(l Line) {
		lines["one"] = append(lines["one"], l)
	})
	lw.Add("two", func(l Line) {
		lines["two"] = append(lines["two"], l)
	})

	lw.Handle(Line{Line: []byte("test1")})
	lw.Handle(Line{Line: []byte("test2")})
	assert.Len(lines["one"], 2)
	assert.Len(lines["two"], 2)

	lw.Remove("two")

	lw.Handle(Line{Line: []byte("test2")})
	assert.Len(lines["one"], 3)
	assert.Len(lines["two"], 2)

	removeTest := new(LineHandlers)
	removeTest.Remove("something")
}
