package web

import (
	"fmt"
	"testing"
	"time"
	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestNestMiddleware(t *testing.T) {
	assert := assert.New(t)

	var callIndex int

	assert.NotNil(NestMiddleware(func(_ *Ctx) Result { return nil }))

	var mw1Called int
	mw1 := func(action Action) Action {
		return func(ctx *Ctx) Result {
			mw1Called = callIndex
			callIndex = callIndex + 1
			return action(ctx)
		}
	}

	var mw2Called int
	mw2 := func(action Action) Action {
		return func(ctx *Ctx) Result {
			mw2Called = callIndex
			callIndex = callIndex + 1
			return action(ctx)
		}
	}

	var mw3Called int
	mw3 := func(action Action) Action {
		return func(ctx *Ctx) Result {
			mw3Called = callIndex
			callIndex = callIndex + 1
			return action(ctx)
		}
	}

	nested := NestMiddleware(func(ctx *Ctx) Result { return nil }, mw2, mw3, mw1)

	nested(nil)

	assert.Equal(2, mw2Called)
	assert.Equal(1, mw3Called)
	assert.Equal(0, mw1Called)
}

func TestPathRedirectHandler(t *testing.T) {
	assert := assert.New(t)

	redirect := PathRedirectHandler("/foo")

	newURL := redirect(NewCtx(nil, webutil.NewMockRequest("GET", "/notfoo")))
	assert.Equal("/foo", newURL.Path)
}

func TestBoolValue(t *testing.T) {
	assert := assert.New(t)

	boolValue, err := BoolValue("true", fmt.Errorf("test"))
	assert.Equal(fmt.Errorf("test"), err)
	assert.False(boolValue)

	boolValue, err = BoolValue("not-bool", nil)
	assert.Equal(fmt.Errorf("invalid boolean value"), err)
	assert.False(boolValue)

	boolValue, err = BoolValue("true", nil)
	assert.Nil(err)
	assert.True(boolValue)

	boolValue, err = BoolValue("1", nil)
	assert.Nil(err)
	assert.True(boolValue)

	boolValue, err = BoolValue("yes", nil)
	assert.Nil(err)
	assert.True(boolValue)

	boolValue, err = BoolValue("0", nil)
	assert.Nil(err)
	assert.False(boolValue)

	boolValue, err = BoolValue("false", nil)
	assert.Nil(err)
	assert.False(boolValue)

	boolValue, err = BoolValue("no", nil)
	assert.Nil(err)
	assert.False(boolValue)
}

func TestIntValue(t *testing.T) {
	assert := assert.New(t)

	value, err := IntValue("1", fmt.Errorf("test"))
	assert.Equal(fmt.Errorf("test"), err)
	assert.Zero(value)

	value, err = IntValue("kdjaflsdf", nil)
	assert.NotNil(err)
	assert.Zero(value)

	value, err = IntValue("1234", nil)
	assert.Nil(err)
	assert.Equal(1234, value)
}

func TestInt64Value(t *testing.T) {
	assert := assert.New(t)

	value, err := Int64Value("1", fmt.Errorf("test"))
	assert.Equal(fmt.Errorf("test"), err)
	assert.Zero(value)

	value, err = Int64Value("kdjaflsdf", nil)
	assert.NotNil(err)
	assert.Zero(value)

	value, err = Int64Value("1234", nil)
	assert.Nil(err)
	assert.Equal(1234, value)
}

func TestFloat64Value(t *testing.T) {
	assert := assert.New(t)

	value, err := Float64Value("1", fmt.Errorf("test"))
	assert.Equal(fmt.Errorf("test"), err)
	assert.Zero(value)

	value, err = Float64Value("kdjaflsdf", nil)
	assert.NotNil(err)
	assert.Zero(value)

	value, err = Float64Value("1234.23", nil)
	assert.Nil(err)
	assert.Equal(1234.23, value)
}

func TestDurationValue(t *testing.T) {
	assert := assert.New(t)

	value, err := DurationValue("1", fmt.Errorf("test"))
	assert.Equal(fmt.Errorf("test"), err)
	assert.Zero(value)

	value, err = DurationValue("kdjaflsdf", nil)
	assert.NotNil(err)
	assert.Zero(value)

	value, err = DurationValue("10s", nil)
	assert.Nil(err)
	assert.Equal(10*time.Second, value)
}

func TestStringValue(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("foo", StringValue("foo", fmt.Errorf("not-foo")))
	assert.Equal("foo", StringValue("foo", nil))
}

func TestCSVValue(t *testing.T) {
	assert := assert.New(t)

	value, err := CSVValue("1", fmt.Errorf("test"))
	assert.Equal(fmt.Errorf("test"), err)
	assert.Empty(value)

	value, err = CSVValue("foo,bar,baz", nil)
	assert.Nil(err)
	assert.Len(value, 3)
}

func TestBase64URL(t *testing.T) {
	assert := assert.New(t)
	bs := []byte("hello")
	enc := Base64URLEncode(bs)
	assert.NotEmpty(enc)

	out, err := Base64URLDecode(enc)
	assert.Nil(err)
	assert.Equal(string(bs), string(out))
}

func TestParseInt32(t *testing.T) {
	assert := assert.New(t)
	i := ParseInt32("10")
	assert.Equal(10, i)
	i = ParseInt32("hbd")
	assert.Equal(0, i)
}

func TestNewCookie(t *testing.T) {
	assert := assert.New(t)
	c := NewCookie("hello", "world")
	assert.NotNil(c)
	assert.Equal("hello", c.Name)
	assert.Equal("world", c.Value)
}
