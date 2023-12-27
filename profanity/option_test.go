/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestConfigOptions(t *testing.T) {
	assert := assert.New(t)

	p := &Profanity{}

	assert.False(p.Config.VerboseOrDefault())
	OptVerbose(true)(p)
	assert.True(p.Config.VerboseOrDefault())

	assert.False(p.Config.DebugOrDefault())
	OptDebug(true)(p)
	assert.True(p.Config.DebugOrDefault())

	assert.False(p.Config.ExitFirstOrDefault())
	OptExitFirst(true)(p)
	assert.True(p.Config.ExitFirstOrDefault())

	assert.Empty(p.Config.Root)
	OptRoot("../foo")(p)
	assert.Equal("../foo", p.Config.Root)

	assert.Equal(DefaultRulesFile, p.Config.RulesFileOrDefault())
	OptRulesFile("my_rules.yml")(p)

	assert.Empty(p.Config.Files.Include)
	OptIncludeFiles("foo", "bar", "baz")(p)
	assert.Equal([]string{"foo", "bar", "baz"}, p.Config.Files.Include)

	assert.Empty(p.Config.Files.Exclude)
	OptExcludeFiles("foo", "bar", "baz")(p)
	assert.Equal([]string{"foo", "bar", "baz"}, p.Config.Files.Exclude)
}
