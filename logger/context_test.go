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
	assert.True(GetTimestamp(context.Background()).IsZero())
}

func TestContextWithScopePath(t *testing.T) {
	assert := assert.New(t)

	path := []string{"one", "two"}
	path2 := []string{"two", "three"}
	assert.Equal(path, GetScopePath(WithScopePath(context.Background(), path...)))
	assert.Equal(path, GetScopePath(WithScopePath(WithScopePath(context.Background(), path2...), path...)))
	assert.Nil(GetScopePath(context.Background()))
}

func TestContextWithLabels(t *testing.T) {
	assert := assert.New(t)

	labels := Labels{"one": "two"}
	labels2 := Labels{"two": "three"}
	assert.Equal(labels, GetLabels(WithLabels(context.Background(), labels)))
	assert.Equal(labels, GetLabels(WithLabels(WithLabels(context.Background(), labels2), labels)))
	assert.Nil(GetLabels(context.Background()))
}
