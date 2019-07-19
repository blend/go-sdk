package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

type zeroTest struct {
	ID    int
	Value string
}

func TestZero(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Input    interface{}
		Expected error
	}{
		{
			Input:    0,
			Expected: nil,
		},
		{
			Input:    1,
			Expected: ErrZero,
		},
		{
			Input:    "",
			Expected: nil,
		},
		{
			Input:    "foo",
			Expected: ErrZero,
		},
		{
			Input:    zeroTest{},
			Expected: nil,
		},
		{
			Input:    zeroTest{ID: 2},
			Expected: ErrZero,
		},
	}

	for index, tc := range testCases {
		verr := Zero(tc.Input)()
		assert.Equal(tc.Expected, ex.ErrClass(ex.ErrInner(verr)), index)
	}
}

func TestNil(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Nil(nil)()
	assert.Nil(verr)

	var nilPtr *string
	verr = Nil(nilPtr)()
	assert.Nil(verr)

	verr = Nil("foo")()
	assert.NotNil(verr)
	assert.Equal(ErrNil, ex.ErrInner(verr))
}

func TestNotNil(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = NotNil("foo")()
	assert.Nil(verr)

	verr = NotNil(nil)()
	assert.NotNil(verr)
	assert.Equal(ErrNotNil, ex.ErrInner(verr))

	var nilPtr *string
	verr = NotNil(nilPtr)()
	assert.NotNil(verr)
	assert.Equal(ErrNotNil, ex.ErrInner(verr))
}

func TestEquals(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Equals("foo")("foo")()
	assert.Nil(verr)

	verr = Equals(nil)(nil)()
	assert.Nil(verr)

	verr = Equals("foo")("bar")()
	assert.NotNil(verr)
	assert.Equal(ErrEquals, ex.ErrInner(verr))

	verr = Equals(nil)("foo")()
	assert.NotNil(verr)
	assert.Equal(ErrEquals, ex.ErrInner(verr))
}

func TestNotEquals(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = NotEquals("foo")("bar")()
	assert.Nil(verr)

	verr = NotEquals(nil)("foo")()
	assert.Nil(verr)

	verr = NotEquals("foo")("foo")()
	assert.NotNil(verr)
	assert.Equal(ErrNotEquals, ex.ErrInner(verr))

	verr = NotEquals(nil)(nil)()
	assert.NotNil(verr)
	assert.Equal(ErrNotEquals, ex.ErrInner(verr))
}

func TestAllow(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Allow("foo")("foo", "bar", "baz")()
	assert.Nil(verr)
	verr = Allow("bar")("foo", "bar", "baz")()
	assert.Nil(verr)
	verr = Allow("baz")("foo", "bar", "baz")()
	assert.Nil(verr)

	verr = Allow("what")("foo", "bar", "baz")()
	assert.NotNil(verr)
	assert.Equal(ErrAllowed, ex.ErrInner(verr))
}

func TestDisallow(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Disallow("what")("foo", "bar", "baz")()
	assert.Nil(verr)

	verr = Disallow("foo")("foo", "bar", "baz")()
	assert.NotNil(verr)
	assert.Equal(ErrDisallowed, ex.ErrInner(verr))
	verr = Disallow("bar")("foo", "bar", "baz")()
	assert.NotNil(verr)
	assert.Equal(ErrDisallowed, ex.ErrInner(verr))
	verr = Disallow("baz")("foo", "bar", "baz")()
	assert.NotNil(verr)
	assert.Equal(ErrDisallowed, ex.ErrInner(verr))
}
