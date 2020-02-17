package main

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/profanity"
)

func TestConfigExampleParses(t *testing.T) {
	assert := assert.New(t)

	p := new(profanity.Profanity)

	rules, err := p.RulesFromReader("config.yml", bytes.NewReader([]byte(configExample)))
	assert.Nil(err)
	assert.Len(rules, 3)
}
