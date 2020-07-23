package webutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestIsValidMethod(t *testing.T) {
	assert := assert.New(t)

	methods := []string{
		MethodGet,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
		MethodOptions,
	}

	for _, method := range methods {
		assert.True(IsValidMethod(method))
	}

	assert.False(IsValidMethod("\n"))
}
