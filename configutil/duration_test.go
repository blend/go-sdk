package configutil

import (
	"testing"
	"time"
	"github.com/blend/go-sdk/assert"
)

func TestDuration(t *testing.T) {
	assert := assert.New(t)

	d := Duration(time.Second)

	ret, err := d.Duration()
	assert.Nil(err)
	assert.NotNil(ret)
	assert.Equal(time.Second, *ret)
}
