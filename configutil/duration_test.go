package configutil

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestDuration(t *testing.T) {
	assert := assert.New(t)

	d := Duration(0)
	ptr, err := d.Duration(context.TODO())
	assert.Nil(ptr)
	assert.Nil(err)

	d = Duration(time.Second)
	ptr, err = d.Duration(context.TODO())
	assert.Nil(err)
	assert.NotNil(ptr)
	assert.Equal(time.Second, *ptr)
}
