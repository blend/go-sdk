package webutil

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMustParseURL(t *testing.T) {
	assert := assert.New(t)

	output := MustParseURL("https://foo.bar.com/example-string?buzz=muzz")
	assert.NotNil(output)
	assert.Equal("https", output.Scheme)
	assert.Equal("foo.bar.com", output.Host)
	assert.Equal("/example-string", output.Path)
	assert.NotEmpty(output.Query())

	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
			}
		}()
		output = MustParseURL(":not-a-url-at-all")
	}()
	assert.NotNil(err)
}

func TestURLWithScheme(t *testing.T) {
	assert := assert.New(t)

	original := MustParseURL("https://foo.bar.com/example-string?buzz=muzz")
	assert.Equal("http", URLWithScheme(original, "http").Scheme)
}

func TestURLWithHost(t *testing.T) {
	assert := assert.New(t)

	original := MustParseURL("https://foo.bar.com/example-string?buzz=muzz")
	assert.Equal("blend.com", URLWithHost(original, "blend.com").Host)
}

func TestURLWithPort(t *testing.T) {
	assert := assert.New(t)

	original := MustParseURL("https://foo.bar.com/example-string?buzz=muzz")
	assert.Equal("foo.bar.com:8443", URLWithPort(original, "8443").Host)
	assert.Equal("8443", URLWithPort(original, "8443").Port())

	originalWithPort := MustParseURL("https://foo.bar.com:5000/example-string?buzz=muzz")
	assert.Equal("foo.bar.com:5001", URLWithPort(originalWithPort, "5001").Host)
	assert.Equal("5001", URLWithPort(originalWithPort, "5001").Port())
}

func TestURLWithPath(t *testing.T) {
	assert := assert.New(t)

	original := MustParseURL("https://foo.bar.com/example-string?buzz=muzz")
	assert.Equal("not-example-string", URLWithPath(original, "not-example-string").Path)
}

func TestURLWithRawQuery(t *testing.T) {
	assert := assert.New(t)

	original := MustParseURL("https://foo.bar.com/example-string?buzz=muzz")
	assert.Equal("dog=cool", URLWithRawQuery(original, "dog=cool").RawQuery)
}

func TestURLWithQuery(t *testing.T) {
	assert := assert.New(t)

	original := MustParseURL("https://foo.bar.com/example-string?buzz=muzz")
	assert.Equal("buzz=muzz", original.RawQuery)
	assert.Equal("buzz=muzz&dog=cool", URLWithQuery(original, "dog", "cool").RawQuery)
}
