package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

// MaybeValidated is a dummy type.
type MaybeValidated struct{}

// Validate implements Validated.
func (mv MaybeValidated) Validate() error { return nil }

func TestIsValidated(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsValidated(MaybeValidated{}))
	assert.False(IsValidated("NOPE"))
}

func TestAsValidated(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(AsValidated(MaybeValidated{}).Validate())
}
