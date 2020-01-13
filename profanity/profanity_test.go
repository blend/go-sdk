package profanity

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestProfanityRulesFromPath(t *testing.T) {
	assert := assert.New(t)

	profanity := &Profanity{}

	rules, err := profanity.RulesFromPath("../" + DefaultRulesFile)
	assert.Nil(err)
	assert.NotEmpty(rules)
}

func TestProfanityReadRules(t *testing.T) {
	assert := assert.New(t)

	profanity := &Profanity{
		Config: Config{
			RulesFile: DefaultRulesFile,
		},
	}

	rules, err := profanity.ReadRules("../")
	assert.Nil(err)
	assert.NotEmpty(rules)
}
