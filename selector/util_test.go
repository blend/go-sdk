package selector

import (
	"fmt"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCheckKey(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(CheckKey("foo"))
	assert.Nil(CheckKey("bar/foo"))
	assert.Nil(CheckKey("bar.io/foo"))
	assert.NotNil(CheckKey("_foo"))
	assert.NotNil(CheckKey("-foo"))
	assert.NotNil(CheckKey("foo-"))
	assert.NotNil(CheckKey("foo_"))
	assert.NotNil(CheckKey("bar/foo/baz"))

	assert.NotNil(CheckKey(""), "should error on empty keys")

	assert.NotNil(CheckKey("/foo"), "should error on empty dns prefixes")
	superLongDNSPrefixed := fmt.Sprintf("%s/%s", strings.Repeat("a", MaxDNSPrefixLen), strings.Repeat("a", MaxKeyLen))
	assert.Nil(CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
	superLongDNSPrefixed = fmt.Sprintf("%s/%s", strings.Repeat("a", MaxDNSPrefixLen+1), strings.Repeat("a", MaxKeyLen))
	assert.NotNil(CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
	superLongDNSPrefixed = fmt.Sprintf("%s/%s", strings.Repeat("a", MaxDNSPrefixLen+1), strings.Repeat("a", MaxKeyLen+1))
	assert.NotNil(CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
	superLongDNSPrefixed = fmt.Sprintf("%s/%s", strings.Repeat("a", MaxDNSPrefixLen), strings.Repeat("a", MaxKeyLen+1))
	assert.NotNil(CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
}

func TestCheckKeyK8S(t *testing.T) {
	assert := assert.New(t)

	values := []string{
		// the "good" cases
		"simple",
		"now-with-dashes",
		"1-starts-with-num",
		"1234",
		"simple/simple",
		"now-with-dashes/simple",
		"now-with-dashes/now-with-dashes",
		"now.with.dots/simple",
		"now-with.dashes-and.dots/simple",
		"1-num.2-num/3-num",
		"1234/5678",
		"1.2.3.4/5678",
		"Uppercase_Is_OK_123",
		"example.com/Uppercase_Is_OK_123",
		"requests.storage-foo",
		strings.Repeat("a", 63),
		strings.Repeat("a", 253) + "/" + strings.Repeat("b", 63),
	}
	badValues := []string{
		// the "bad" cases
		"nospecialchars%^=@",
		"cantendwithadash-",
		"-cantstartwithadash-",
		"only/one/slash",
		"Example.com/abc",
		"example_com/abc",
		"example.com/",
		"/simple",
		strings.Repeat("a", 64),
		strings.Repeat("a", 254) + "/abc",
	}
	for _, val := range values {
		assert.Nil(CheckKey(val))
	}
	for _, val := range badValues {
		assert.NotNil(CheckKey(val))
	}
}

func TestCheckValue(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(CheckValue(""), "should not error on empty values")
	assert.Nil(CheckValue("foo"))
	assert.Nil(CheckValue("bar_baz"))
	assert.NotNil(CheckValue("_bar_baz"))
	assert.NotNil(CheckValue("bar_baz_"))
	assert.NotNil(CheckValue("_bar_baz_"))
}

func TestIsAlpha(t *testing.T) {
	assert := assert.New(t)

	assert.True(isAlpha('A'))
	assert.True(isAlpha('a'))
	assert.True(isAlpha('Z'))
	assert.True(isAlpha('z'))
	assert.True(isAlpha('0'))
	assert.True(isAlpha('9'))
	assert.True(isAlpha('함'))
	assert.True(isAlpha('é'))
	assert.False(isAlpha('-'))
	assert.False(isAlpha('/'))
	assert.False(isAlpha('~'))
}

func TestCheckLabels(t *testing.T) {
	assert := assert.New(t)

	goodLabels := Labels{"foo": "bar", "foo.com/bar": "baz"}
	assert.Nil(CheckLabels(goodLabels))
	badLabels := Labels{"foo": "bar", "_foo.com/bar": "baz"}
	assert.NotNil(CheckLabels(badLabels))
}

func TestIsPermitted(t *testing.T) {
	assert := assert.New(t)

	assert.False(isPermitted('%', nil))
	assert.False(isPermitted('%', map[rune]bool{}))
	assert.False(isPermitted('%', map[rune]bool{'*': true}))
	assert.True(isPermitted('%', map[rune]bool{'*': true, '%': true}))
}
