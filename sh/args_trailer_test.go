package sh

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestParseFlagsTrailer(t *testing.T) {
	assert := assert.New(t)

	parsed, err := ArgsTrailer("foo", "bar")
	assert.True(ex.Is(err, ErrFlagsNoTrailer))
	assert.Empty(parsed)

	parsed, err = ArgsTrailer("foo", "bar", "--")
	assert.True(ex.Is(err, ErrFlagsNoTrailer))
	assert.Empty(parsed)

	parsed, err = ArgsTrailer("foo", "bar", "--", "echo", "'things'")
	assert.Nil(err)
	assert.Equal([]string{"echo", "'things'"}, parsed)
}
