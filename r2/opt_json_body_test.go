package r2

import (
	"io/ioutil"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestOptJSONBody(t *testing.T) {
	assert := assert.New(t)

	object := map[string]interface{}{"foo": "bar"}

	opt := OptJSONBody(object)

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Body)

	contents, err := ioutil.ReadAll(req.Body)
	assert.Nil(err)
	assert.Equal(`{"foo":"bar"}`, string(contents))

	assert.Equal(ContentTypeApplicationJSON, req.Header.Get("Content-Type"))
}
