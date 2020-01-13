package yaml_test

import (
	"strconv"
	"testing"
	"github.com/blend/go-sdk/assert"
)

// TestMain is the testing entrypoint.
func TestMain(m *testing.M) {
	assert.Main(m)
}

// MustUnquote unquotes a string, panicing if there is an issue.
func MustUnquote(str string) string {
	value, err := strconv.Unquote(str)
	if err != nil {
		panic(err)
	}
	return value
}
