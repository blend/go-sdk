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
			Input:    int(0),
			Expected: nil,
		},
		{
			Input:    int(1),
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
		assert.Equal(tc.Expected, ex.ErrInner(Zero(tc.Input)), index)
	}
}

func TestNil(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Nil(nil)
	assert.Nil(verr)

	var nilPtr *string
	verr = Nil(nilPtr)
	assert.Nil(verr)

	verr = Nil("foo")
	assert.NotNil(verr)
	assert.Equal(ErrNil, ex.ErrInner(verr))
}

func TestNotNil(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = NotNil("foo")
	assert.Nil(verr)

	verr = NotNil(nil)
	assert.NotNil(verr)
	assert.Equal(ErrNotNil, ex.ErrInner(verr))

	var nilPtr *string
	verr = NotNil(nilPtr)
	assert.NotNil(verr)
	assert.Equal(ErrNotNil, ex.ErrInner(verr))
}

func TestEquals(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Equals("foo")("foo")
	assert.Nil(verr)

	verr = Equals(nil)(nil)
	assert.Nil(verr)

	verr = Equals("foo")("bar")
	assert.NotNil(verr)
	assert.Equal(ErrEquals, ex.ErrInner(verr))

	verr = Equals(nil)("foo")
	assert.NotNil(verr)
	assert.Equal(ErrEquals, ex.ErrInner(verr))
}

func TestNotEquals(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = NotEquals("foo")("bar")
	assert.Nil(verr)

	verr = NotEquals(nil)("foo")
	assert.Nil(verr)

	verr = NotEquals("foo")("foo")
	assert.NotNil(verr)
	assert.Equal(ErrNotEquals, ex.ErrInner(verr))

	verr = NotEquals(nil)(nil)
	assert.NotNil(verr)
	assert.Equal(ErrNotEquals, ex.ErrInner(verr))
}

func TestAllow(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Allow("foo", "bar", "baz")("foo")
	assert.Nil(verr)
	verr = Allow("foo", "bar", "baz")("bar")
	assert.Nil(verr)
	verr = Allow("foo", "bar", "baz")("baz")
	assert.Nil(verr)

	verr = Allow("foo", "bar", "baz")("what")
	assert.NotNil(verr)
	assert.Equal(ErrAllowed, ex.ErrInner(verr))
}

func TestDisallow(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Disallow("foo", "bar", "baz")("what")
	assert.Nil(verr)

	verr = Disallow("foo", "bar", "baz")("foo")
	assert.NotNil(verr)
	assert.Equal(ErrDisallowed, ex.ErrInner(verr))
	verr = Disallow("foo", "bar", "baz")("bar")
	assert.NotNil(verr)
	assert.Equal(ErrDisallowed, ex.ErrInner(verr))
	verr = Disallow("foo", "bar", "baz")("baz")
	assert.NotNil(verr)
	assert.Equal(ErrDisallowed, ex.ErrInner(verr))
}
