/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
