package webutil

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestSplitColon(t *testing.T) {
	assert := assert.New(t)

	// Missing ":"
	input := "some text"
	_, _, err := SplitColon(input)
	assert.True(ErrIsInvalidSplitColonInput(err))
	assert.Equal(fmt.Sprintf(`input: %q`, input), ex.ErrMessage(err))

	// No text before the ":"
	input = ":p4ssw0rd"
	_, _, err = SplitColon(input)
	assert.True(ErrIsInvalidSplitColonInput(err))
	assert.Equal(fmt.Sprintf(`input: %q`, input), ex.ErrMessage(err))

	// No text after the ":"
	input = "user@mail.invalid:"
	_, _, err = SplitColon(input)
	assert.True(ErrIsInvalidSplitColonInput(err))
	assert.Equal(fmt.Sprintf(`input: %q`, input), ex.ErrMessage(err))

	// Valid input value
	var first, second string
	first, second, err = SplitColon("cake:eat-it-too")
	assert.Nil(err)
	assert.Equal(first, "cake")
	assert.Equal(second, "eat-it-too")
}
