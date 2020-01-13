package r2

import (
	"net/url"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestOptPostForm(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptPostForm(url.Values{"bar": []string{"baz, buzz"}}))
	assert.NotNil(r.PostForm)
	assert.NotEmpty(r.PostForm.Get("bar"))
}

func TestOptPostFormValue(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptPostFormValue("bar", "baz"))
	assert.NotNil(r.PostForm)
	assert.Equal("baz", r.PostForm.Get("bar"))
}
