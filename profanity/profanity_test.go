package profanity

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestProfanityRulesFromPath(t *testing.T) {
	assert := assert.New(t)

	profanity := &Profanity{}

	rules, err := profanity.RulesFromPath("testdata/rules.yml")
	assert.Nil(err)
	assert.NotEmpty(rules)
}

func TestProfanityReadRules(t *testing.T) {
	assert := assert.New(t)

	profanity := &Profanity{
		Config: Config{
			RulesFile: "rules.yml",
		},
	}

	rules, err := profanity.ReadRules("testdata/")
	assert.Nil(err)
	assert.NotEmpty(rules)
}
