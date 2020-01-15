package sentry

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestContextFingerprint(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(GetFingerprint(nil))
	assert.Nil(GetFingerprint(context.Background()))
	assert.Nil(GetFingerprint(context.WithValue(context.Background(), contextFingerprintKey{}, 1234)))

	assert.Equal([]string{"foo", "bar"}, GetFingerprint(WithFingerprint(context.Background(), "foo", "bar")))
}
