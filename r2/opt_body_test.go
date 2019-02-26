package r2

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptBody(t *testing.T) {
	assert := assert.New(t)

	req := New("https://foo.bar.local")

	assert.Nil(OptBody(ioutil.NopCloser(bytes.NewBufferString("this is only a test")))(req))
	assert.NotNil(req.Body)

	contents, err := ioutil.ReadAll(req.Body)
	assert.Nil(err)
	assert.Equal("this is only a test", string(contents))
}
