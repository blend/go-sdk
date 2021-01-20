package protoutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/protoutil/testdata"
)

func Test_MessageTypeName(t *testing.T) {
	its := assert.New(t)

	its.Equal("testdata.Message", MessageTypeName(new(testdata.Message)))
}
