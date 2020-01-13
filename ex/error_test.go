package ex

import (
	"encoding/json"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestClassMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	err := Class("this is only a test")
	contents, marshalErr := json.Marshal(err)
	assert.Nil(marshalErr)
	assert.Equal(`"this is only a test"`, string(contents))
}
