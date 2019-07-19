package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestMapKeys(t *testing.T) {
	assert := assert.New(t)

	var err error
	err = Map("foo").Keys("a", "b", "c")()
	assert.NotNil(err)
	assert.Equal(ErrInstanceNotMap, ex.ErrClass(err))

	bag := map[string]int{"foo": 1, "bar": 2, "baz": 3}
	var verr error
	verr = Map(bag).Keys("foo", "baz")()
	assert.Nil(verr)

	verr = Map(bag).Keys("foo", "buzz")()
	assert.NotNil(verr)
	assert.Equal(ErrMapKeys, ex.ErrInner(verr))
}
