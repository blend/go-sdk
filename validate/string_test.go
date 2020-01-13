package validate

import (
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/uuid"
)

func TestStringRequired(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String(nil).Required()()
	assert.NotNil(verr)
	assert.Equal(ErrStringRequired, ErrCause(ErrStringRequired))

	bad := ""
	verr = String(&bad).Required()()
	assert.NotNil(verr)
	assert.Equal(ErrStringRequired, ErrCause(ErrStringRequired))

	good := "ok!"
	verr = String(&good).Required()()
	assert.Nil(verr)
}

func TestStringMin(t *testing.T) {
	assert := assert.New(t)

	var verr error
	bad := "large"
	verr = String(&bad).MinLen(3)()
	assert.Nil(verr)

	verr = String(nil).MinLen(3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, ErrCause(verr))

	good := "a"
	verr = String(&good).MinLen(3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, ErrCause(verr))
}

func TestStringMaxlen(t *testing.T) {
	assert := assert.New(t)

	var verr error
	bad := "a"
	verr = String(&bad).MaxLen(3)()
	assert.Nil(verr)

	verr = String(nil).MaxLen(3)()
	assert.Nil(verr)

	good := "large"
	verr = String(&good).MaxLen(3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMax, ErrCause(verr))
}

func TestStringBetweenLen(t *testing.T) {
	assert := assert.New(t)

	var verr error
	good := "ok"
	verr = String(&good).BetweenLen(1, 3)()
	assert.Nil(verr)

	bad := "large"
	verr = String(&bad).BetweenLen(1, 3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMax, ErrCause(verr))

	verr = String(nil).BetweenLen(2, 5)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, ErrCause(verr))

	bad = "a"
	verr = String(&bad).BetweenLen(2, 5)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, ErrCause(verr))
}

func TestStringMatches(t *testing.T) {
	assert := assert.New(t)

	var verr error
	good := "a foo"
	verr = String(&good).Matches("foo$")()
	assert.Nil(verr)

	verr = String(nil).Matches("foo$")()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrStringMatches, ErrCause(verr))

	bad := "foo not"
	verr = String(&bad).Matches("foo$")()
	assert.NotNil(verr)
	assert.Equal(ErrStringMatches, ErrCause(verr))
}

func TestStringMatchesError(t *testing.T) {
	assert := assert.New(t)

	var err error
	good := "a foo"
	err = String(&good).Matches("((")() // this should be an invalid regex "(("
	assert.NotNil(err)
	assert.NotEqual(ErrValidation, ex.ErrClass(err))
}

func TestStringIsUpper(t *testing.T) {
	assert := assert.New(t)

	var verr error
	good := "FOO"
	verr = String(&good).IsUpper()()
	assert.Nil(verr)

	verr = String(nil).IsUpper()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrStringIsUpper, ErrCause(verr))

	bad := "FOo"
	verr = String(&bad).IsUpper()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsUpper, ErrCause(verr))
}

func TestStringIsLower(t *testing.T) {
	assert := assert.New(t)

	var verr error
	good := "foo"
	verr = String(&good).IsLower()()
	assert.Nil(verr)

	verr = String(nil).IsLower()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrStringIsLower, ErrCause(verr))

	bad := "foO"
	verr = String(&bad).IsLower()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsLower, ErrCause(verr))
}

func TestStringIsTitle(t *testing.T) {
	assert := assert.New(t)

	var verr error
	good := strings.ToTitle("this is a test")
	verr = String(&good).IsTitle()()
	assert.Nil(verr)

	verr = String(nil).IsTitle()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrStringIsTitle, ErrCause(verr))

	bad := "this is a test"
	verr = String(&bad).IsTitle()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsTitle, ErrCause(verr))
}

func TestStringIsUUID(t *testing.T) {
	assert := assert.New(t)

	var verr error
	good := uuid.V4().String()
	verr = String(&good).IsUUID()()
	assert.Nil(verr)

	verr = String(nil).IsUUID()()
	assert.NotNil(verr)
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrStringIsUUID, ErrCause(verr))

	good = uuid.V4().ToFullString()
	verr = String(&good).IsUUID()()
	assert.Nil(verr)

	bad := "asldkfjaslkfjasdlfa"
	verr = String(&bad).IsUUID()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsUUID, ErrCause(verr))
}

func TestStringIsEmail(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String(nil).IsEmail()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrStringIsEmail, ErrCause(verr))

	good := "foo@bar.com"
	verr = String(&good).IsEmail()()
	assert.Nil(verr)

	good = "foo@bar"
	verr = String(&good).IsEmail()()
	assert.Nil(verr)

	good = "foo+foo@bar.com"
	verr = String(&good).IsEmail()()
	assert.Nil(verr)

	bad := "this is a test"
	verr = String(&bad).IsEmail()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsEmail, ErrCause(verr))
}

func TestStringIsURI(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String(nil).IsURI()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrStringIsURI, ErrCause(verr))

	good := "https://foo.com"

	verr = String(&good).IsURI()()
	assert.Nil(verr)

	bad := "this is a test"
	verr = String(&bad).IsURI()()
	assert.NotNil(verr)
	assert.Equal(bad, ErrValue(verr))
	assert.Equal(ErrStringIsURI, ErrCause(verr))
}

func TestStringIsIP(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String(nil).IsIP()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrStringIsIP, ErrCause(verr))

	good := "127.0.0.1"
	verr = String(&good).IsIP()()
	assert.Nil(verr)
	good = "172.217.0.46"
	verr = String(&good).IsIP()()
	assert.Nil(verr)
	good = "2607:f8b0:4005:802::200e"
	verr = String(&good).IsIP()()
	assert.Nil(verr)
	good = "::1"
	verr = String(&good).IsIP()()
	assert.Nil(verr)

	bad := ""
	verr = String(&bad).IsIP()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsIP, ErrCause(verr))

	bad = "this is a test"
	verr = String(&bad).IsIP()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsIP, ErrCause(verr))
}

func TestStringIsSlug(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String(nil).IsSlug()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrStringIsSlug, ErrCause(verr))

	good := "abcdefghijklmnopqrstuvwxyz"
	verr = String(&good).IsSlug()()
	assert.Nil(verr)

	good = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	verr = String(&good).IsSlug()()
	assert.Nil(verr)

	good = "0123456789"
	verr = String(&good).IsSlug()()
	assert.Nil(verr)

	good = "_-"
	verr = String(&good).IsSlug()()
	assert.Nil(verr)

	good = "shortcut_service"
	verr = String(&good).IsSlug()()
	assert.Nil(verr)

	good = "Shortcut-Service"
	verr = String(&good).IsSlug()()
	assert.Nil(verr)

	bad := "this/../is/../hacking?"
	verr = String(&bad).IsSlug()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsSlug, ErrCause(verr))
}

func TestStringIsOneOf(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String(nil).IsOneOf()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsOneOf, ErrCause(verr))
	assert.Nil(ErrValue(verr))
	assert.NotNil(ErrMessage(verr))

	good := "foo"
	verr = String(&good).IsOneOf("foo", "bar")()
	assert.Nil(verr)

	bad := "bad"
	verr = String(&bad).IsOneOf("foo", "bar")()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrStringIsOneOf, ErrCause(verr))
	assert.Equal("foo, bar", ErrMessage(verr))
}
