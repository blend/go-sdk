package selector

import (
	"testing"

	assert "github.com/blend/go-sdk/assert"
)

func TestSkipValidation(t *testing.T) {
	assert := assert.New(t)

	// Test: star is not a valid value
	_, err := Parse("foo=ba*r")
	assert.Equal(ErrKeyInvalidCharacter, err)
	// Test: SkipValidation should cause invalid values to be ignored
	_, err = Parse("foo=ba*r", SkipValidation)
	assert.Nil(err)
}
func TestOptExtraAlphas(t *testing.T) {
	assert := assert.New(t)

	// Test: =* should result in an invalid op error
	_, err := Parse("foo=*bar")
	assert.Equal(ErrInvalidOperator, err)

	// Test: Permitting '*' should cause the parser to interpret '*' as part of the value
	// and not part of the operator
	_, err = Parse("foo=*bar", OptExtraAlphas('*'))
	assert.Equal(ErrKeyInvalidCharacter, err)

	// Test: Permitting '*' and skipping validation will cause '*'-prefixed values to be allowed
	_, err = Parse("foo=*bar", SkipValidation, OptExtraAlphas('*'))
	assert.Nil(err)
}
