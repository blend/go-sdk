package webutil

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func Test_mustNil(t *testing.T) {
	assert := assert.New(t)

	err := ex.New("bad-bad")
	panicFn := func() { mustNil(err) }
	assert.PanicEqual(err, panicFn)
}

func TestMustParseSameSite(t *testing.T) {
	assert := assert.New(t)

	ss := MustParseSameSite(SameSiteStrict)
	assert.Equal(http.SameSiteStrictMode, ss)
}

func TestParseSameSite(t *testing.T) {
	assert := assert.New(t)

	ss, err := ParseSameSite(SameSiteStrict)
	assert.Nil(err)
	assert.Equal(http.SameSiteStrictMode, ss)

	ss, err = ParseSameSite(SameSiteLax)
	assert.Nil(err)
	assert.Equal(http.SameSiteLaxMode, ss)

	ss, err = ParseSameSite(SameSiteDefault)
	assert.Nil(err)
	assert.Equal(http.SameSiteDefaultMode, ss)

	ss, err = ParseSameSite("not-included")
	assert.Equal(-1, ss)
	assert.NotNil(err)
	assert.True(ex.Is(ErrInvalidSameSite, err))
}
