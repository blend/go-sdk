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
		verr := Any(tc.Input).Zero()
		assert.Equal(tc.Expected, ex.ErrInner(verr), index)
	}
}

func TestNil(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any(nil).Nil()
	assert.Nil(verr)

	var nilPtr *string
	verr = Any(nilPtr).Nil()
	assert.Nil(verr)

	verr = Any("foo").Nil()
	assert.NotNil(verr)
	assert.Equal(ErrNil, ex.ErrInner(verr))
}

func TestNotNil(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any("foo").NotNil()
	assert.Nil(verr)

	verr = Any(nil).NotNil()
	assert.NotNil(verr)
	assert.Equal(ErrNotNil, ex.ErrInner(verr))

	var nilPtr *string
	verr = Any(nilPtr).NotNil()
	assert.NotNil(verr)
	assert.Equal(ErrNotNil, ex.ErrInner(verr))
}

func TestEquals(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any("foo").Equals("foo")()
	assert.Nil(verr)

	verr = Any(nil).Equals(nil)()
	assert.Nil(verr)

	verr = Any("foo").Equals("bar")()
	assert.NotNil(verr)
	assert.Equal(ErrEquals, ex.ErrInner(verr))

	verr = Any(nil).Equals("foo")()
	assert.NotNil(verr)
	assert.Equal(ErrEquals, ex.ErrInner(verr))
}

func TestNotEquals(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any("foo").NotEquals("bar")()
	assert.Nil(verr)

	verr = Any(nil).NotEquals("foo")()
	assert.Nil(verr)

	verr = Any("foo").NotEquals("foo")()
	assert.NotNil(verr)
	assert.Equal(ErrNotEquals, ex.ErrInner(verr))

	verr = Any(nil).NotEquals(nil)()
	assert.NotNil(verr)
	assert.Equal(ErrNotEquals, ex.ErrInner(verr))
}

func TestAllow(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any("foo").Allow("foo", "bar", "baz")()
	assert.Nil(verr)
	verr = Any("bar").Allow("foo", "bar", "baz")()
	assert.Nil(verr)
	verr = Any("baz").Allow("foo", "bar", "baz")()
	assert.Nil(verr)

	verr = Any("what").Allow("foo", "bar", "baz")()
	assert.NotNil(verr)
	assert.Equal(ErrAllowed, ex.ErrInner(verr))
}

func TestDisallow(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any("what").Disallow("foo", "bar", "baz")()
	assert.Nil(verr)

	verr = Any("foo").Disallow("foo", "bar", "baz")()
	assert.NotNil(verr)
	assert.Equal(ErrDisallowed, ex.ErrInner(verr))
	verr = Any("bar").Disallow("foo", "bar", "baz")()
	assert.NotNil(verr)
	assert.Equal(ErrDisallowed, ex.ErrInner(verr))
	verr = Any("baz").Disallow("foo", "bar", "baz")()
	assert.NotNil(verr)
	assert.Equal(ErrDisallowed, ex.ErrInner(verr))
}
