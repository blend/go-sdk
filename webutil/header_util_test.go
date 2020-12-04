package webutil

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestHeaderAny(t *testing.T) {
	assert := assert.New(t)

	assert.True(HeaderAny(http.Header{"Foo": []string{"bar"}}, "foo", "bar"))
	assert.True(HeaderAny(http.Header{"fuzz": []string{"buzz"}, "Foo": []string{"bar"}}, "foo", "bar"))
	assert.False(HeaderAny(http.Header{"fuzz": []string{"buzz"}, "Foo": []string{"bar"}}, "fuzz", "bar"))
	assert.True(HeaderAny(http.Header{"Foo": []string{"bailey,bar"}}, "foo", "bar"))
	assert.True(HeaderAny(http.Header{"Foo": []string{"bar,bailey"}}, "foo", "bar"))
	assert.True(HeaderAny(http.Header{"fuzz": []string{"buzz"}, "Foo": []string{"bar,bailey"}}, "foo", "bar"))
}
