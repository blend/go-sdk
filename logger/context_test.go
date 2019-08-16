package logger

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestContextWithTimestamp(t *testing.T) {
	assert := assert.New(t)

	ts := time.Date(2019, 8, 16, 12, 11, 10, 9, time.UTC)
	assert.Equal(ts, GetTimestamp(WithTimestamp(context.Background(), ts)))
	assert.False(GetTimestamp(context.Background()).IsZero())
}

func TestContextWithScopePath(t *testing.T) {
	assert := assert.New(t)

	path := []string{"one", "two"}
	path2 := []string{"two", "three"}
	assert.Equal(path, GetScopePath(WithScopePath(context.Background(), path...)))
	assert.Equal(path, GetScopePath(WithScopePath(WithScopePath(context.Background(), path2...), path...)))
	assert.Nil(GetScopePath(context.Background()))
}

func TestContextWithFields(t *testing.T) {
	assert := assert.New(t)

	fields := Fields{"one": "two"}
	fields2 := Fields{"two": "three"}
	assert.Equal(fields, GetFields(WithFields(context.Background(), fields)))
	assert.Equal(fields, GetFields(WithFields(WithFields(context.Background(), fields2), fields)))
	assert.Nil(GetFields(context.Background()))
}
