package ansi

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestTableForSlice(t *testing.T) {
	assert := assert.New(t)

	objects := []struct {
		ID   int
		Name string
	}{
		{1, "Foo"},
		{2, "Bar"},
		{3, "Baz"},
	}

	output := new(bytes.Buffer)
	assert.Nil(TableForSlice(output, objects))
	assert.Equal(
		"┌────┬──────┐\n│ ID │ Name │\n├────┼──────┤\n│ 1  │ Foo  │\n│ 2  │ Bar  │\n│ 3  │ Baz  │\n└────┴──────┘\n",
		output.String(),
	)
}
