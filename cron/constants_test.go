package cron

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestErrors(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsJobNotLoaded(ErrJobNotLoaded))
	assert.True(IsJobNotLoaded(ex.New(ErrJobNotLoaded)))
	assert.False(IsJobNotLoaded(ErrJobAlreadyLoaded))
	assert.False(IsJobNotLoaded(ex.New(ErrJobAlreadyLoaded)))

	assert.True(IsJobAlreadyLoaded(ErrJobAlreadyLoaded))
	assert.True(IsJobAlreadyLoaded(ex.New(ErrJobAlreadyLoaded)))
	assert.False(IsJobAlreadyLoaded(ErrJobNotLoaded))
	assert.False(IsJobAlreadyLoaded(ex.New(ErrJobNotLoaded)))

	assert.True(IsJobNotFound(ErrJobNotFound))
	assert.True(IsJobNotFound(ex.New(ErrJobNotFound)))
	assert.False(IsJobNotFound(ErrJobNotLoaded))
	assert.False(IsJobNotFound(ex.New(ErrJobNotLoaded)))
}
