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
	bad := "large"
	verr = String(&bad).MinLen(3)()
	assert.Nil(verr)

	verr = String(nil).MinLen(3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, Cause(verr))

	good := "a"
	verr = String(&good).MinLen(3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, Cause(verr))
}

func TestStringMaxlen(t *testing.T) {
	assert := assert.New(t)

	var verr error
	bad := "a"
	verr = String(&bad).MaxLen(3)()
	assert.Nil(verr)
	good := "large"
	verr = String(&good).MaxLen(3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMax, Cause(verr))
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
	assert.Equal(ErrStringLengthMax, Cause(verr))

	verr = String(nil).BetweenLen(2, 5)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, Cause(verr))

	bad = "a"
	verr = String(&bad).BetweenLen(2, 5)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, Cause(verr))
}

func TestStringMatches(t *testing.T) {
	assert := assert.New(t)

	var verr error
	good := "a foo"
	verr = String(&good).Matches("foo$")()
	assert.Nil(verr)

	bad := "foo not"
	verr = String(&bad).Matches("foo$")()
	assert.NotNil(verr)
	assert.Equal(ErrStringMatches, Cause(verr))
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

	bad := "FOo"
	verr = String(&bad).IsUpper()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsUpper, Cause(verr))
}

func TestStringIsLower(t *testing.T) {
	assert := assert.New(t)

	var verr error
	good := "foo"
	verr = String(&good).IsLower()()
	assert.Nil(verr)

	bad := "foO"
	verr = String(&bad).IsLower()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsLower, Cause(verr))
}

func TestStringIsTitle(t *testing.T) {
	assert := assert.New(t)

	var verr error
	good := strings.ToTitle("this is a test")
	verr = String(&good).IsTitle()()
	assert.Nil(verr)

	bad := "this is a test"
	verr = String(&bad).IsTitle()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsTitle, Cause(verr))
}

func TestStringIsUUID(t *testing.T) {
	assert := assert.New(t)

	var verr error
	good := uuid.V4().String()
	verr = String(&good).IsUUID()()
	assert.Nil(verr)

	good = uuid.V4().ToFullString()
	verr = String(&good).IsUUID()()
	assert.Nil(verr)

	bad := "asldkfjaslkfjasdlfa"
	verr = String(&bad).IsUUID()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsUUID, Cause(verr))
}

func TestStringIsEmail(t *testing.T) {
	assert := assert.New(t)

	var verr error
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
	assert.Equal(ErrStringIsEmail, Cause(verr))
}

func TestStringIsURI(t *testing.T) {
	assert := assert.New(t)

	var verr error
	good := "https://foo.com"
	verr = String(&good).IsURI()()
	assert.Nil(verr)

	bad := "this is a test"
	verr = String(&bad).IsURI()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsURI, Cause(verr))
}

func TestStringIsIP(t *testing.T) {
	assert := assert.New(t)

	var verr error
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
	assert.Equal(ErrStringIsIP, Cause(verr))

	bad = "this is a test"
	verr = String(&bad).IsIP()()
	assert.NotNil(verr)
	assert.Equal(ErrStringIsIP, Cause(verr))
}
