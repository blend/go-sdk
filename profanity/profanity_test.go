/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Profanity_ReadRuleSpecsFile(t *testing.T) {
	assert := assert.New(t)

	profanity := &Profanity{}

	rules, err := profanity.ReadRuleSpecsFile("testdata/rules.yml")
	assert.Nil(err)
	assert.NotEmpty(rules)
}
