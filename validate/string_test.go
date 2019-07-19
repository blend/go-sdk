package validate

import (
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/uuid"
)

func TestStringMin(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String("large").MinLen(3)()
	assert.Nil(verr)
	verr = String("a").MinLen(3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, ex.ErrInner(verr))
}

func TestStringMaxlen(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String("a").MaxLen(3)()
	assert.Nil(verr)
	verr = String("large").MaxLen(3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMax, ex.ErrInner(verr))
}

func TestStringBetweenLen(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String("a").BetweenLen(1, 3)()
	assert.Nil(verr)

	verr = String("large").BetweenLen(1, 3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMax, ex.ErrInner(verr))

	verr = String("").BetweenLen(1, 3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, ex.ErrInner(verr))
}

func TestStringMatches(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String("a foo").Matches("foo$")()
	assert.Nil(verr)

	verr = String("foo not").Matches("foo$")()
	assert.NotNil(verr)
	assert.Equal(ErrStringMatches, ex.ErrInner(verr))
}

func TestStringMatchesError(t *testing.T) {
	assert := assert.New(t)

	var err error
	err = String("a foo").Matches("((")() // this should be an invalid regex "(("
	assert.NotNil(err)
	assert.NotEqual(ErrValidation, ex.ErrClass(err))
}

func TestStringIsUpper(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String("FOO").IsUpper()
	assert.Nil(verr)

	verr = String("FOo").IsUpper()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsUpper, ex.ErrInner(verr))
}

func TestStringIsLower(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String("foo").IsLower()
	assert.Nil(verr)

	verr = String("foO").IsLower()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsLower, ex.ErrInner(verr))
}

func TestStringIsTitle(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String(strings.ToTitle("this is a test")).IsTitle()
	assert.Nil(verr)

	verr = String("this is a test").IsTitle()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsTitle, ex.ErrInner(verr))
}

func TestStringIsUUID(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String(uuid.V4().String()).IsUUID()
	assert.Nil(verr)

	verr = String(uuid.V4().ToFullString()).IsUUID()
	assert.Nil(verr)

	verr = String("asldkfjaslkfjasdlfa").IsUUID()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsUUID, ex.ErrInner(verr))
}

func TestStringIsEmail(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String("foo@bar.com").IsEmail()
	assert.Nil(verr)

	verr = String("foo@bar").IsEmail()
	assert.Nil(verr)

	verr = String("foo+foo@bar.com").IsEmail()
	assert.Nil(verr)

	verr = String("this is a test").IsEmail()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsEmail, ex.ErrInner(verr))
}

func TestStringIsURI(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String("https://foo.com").IsURI()
	assert.Nil(verr)

	verr = String("this is a test").IsURI()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsURI, ex.ErrInner(verr))
}
