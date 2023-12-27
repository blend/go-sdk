/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptQuery(t *testing.T) {
	assert := assert.New(t)

	req := New(TestURL,
		OptQuery(url.Values{
			"huff":	[]string{"buff"},
			"buzz":	[]string{"fuzz"},
		}),
	)
	assert.NotNil(req.Request.URL)
	assert.NotEmpty(req.Request.URL.RawQuery)
	assert.NotEmpty(req.Request.URL.Query())
	assert.Equal("buff", req.Request.URL.Query().Get("huff"))
	assert.Equal("fuzz", req.Request.URL.Query().Get("buzz"))
	assert.Equal("buzz=fuzz&huff=buff", req.Request.URL.RawQuery)
}

func TestOptQueryValue(t *testing.T) {
	assert := assert.New(t)

	req := New(TestURL,
		OptQueryValue("huff", "buff"),
		OptQueryValue("buzz", "fuzz"),
	)
	assert.NotNil(req.Request.URL)
	assert.NotEmpty(req.Request.URL.RawQuery)
	assert.NotEmpty(req.Request.URL.Query())
	assert.Equal("buff", req.Request.URL.Query().Get("huff"))
	assert.Equal("fuzz", req.Request.URL.Query().Get("buzz"))
	assert.Equal("buzz=fuzz&huff=buff&query=value", req.Request.URL.RawQuery)
}

func TestOptQueryValueAdd(t *testing.T) {
	assert := assert.New(t)

	req := New(TestURL,
		OptQueryValue("huff", "buff"),
		OptQueryValueAdd("huff", "fuzz"),
		OptQueryValueAdd("foo", "bar"),
		OptQueryValueAdd("foo", "baz"),
	)
	assert.NotNil(req.Request.URL)
	assert.NotEmpty(req.Request.URL.RawQuery)
	assert.NotEmpty(req.Request.URL.Query())
	assert.Equal([]string{"buff", "fuzz"}, req.Request.URL.Query()["huff"])
	assert.Equal([]string{"bar", "baz"}, req.Request.URL.Query()["foo"])
}
