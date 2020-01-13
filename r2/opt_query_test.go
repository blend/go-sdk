package r2

import (
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptQuery(t *testing.T) {
	assert := assert.New(t)

	req := New("http://foo.bar.com",
		OptQuery(url.Values{
			"huff": []string{"buff"},
			"buzz": []string{"fuzz"},
		}),
	)
	assert.NotNil(req.URL)
	assert.NotEmpty(req.URL.RawQuery)
	assert.NotEmpty(req.URL.Query())
	assert.Equal("buff", req.URL.Query().Get("huff"))
	assert.Equal("fuzz", req.URL.Query().Get("buzz"))
	assert.Equal("buzz=fuzz&huff=buff", req.URL.RawQuery)
}

func TestOptQueryValue(t *testing.T) {
	assert := assert.New(t)

	req := New("http://foo.bar.com",
		OptQueryValue("huff", "buff"),
		OptQueryValue("buzz", "fuzz"),
	)
	assert.NotNil(req.URL)
	assert.NotEmpty(req.URL.RawQuery)
	assert.NotEmpty(req.URL.Query())
	assert.Equal("buff", req.URL.Query().Get("huff"))
	assert.Equal("fuzz", req.URL.Query().Get("buzz"))
	assert.Equal("buzz=fuzz&huff=buff", req.URL.RawQuery)
}
