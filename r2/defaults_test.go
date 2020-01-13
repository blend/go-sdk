package r2

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func optNone() func(r *Request) error {
	return func(r *Request) error {
		return nil
	}
}

func TestDefaults(t *testing.T) {
	assert := assert.New(t)

	def := Defaults([]Option{optNone(), optNone()})
	assert.Len(def, 2)

	def = def.Add(optNone(), optNone())
	assert.Len(def, 4)
}
