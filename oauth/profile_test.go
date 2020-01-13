package oauth

import (
	"testing"
	assert "github.com/blend/go-sdk/assert"
)

func TestProfileUsername(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(Profile{}.Username())
	assert.Equal("foo", Profile{Email: "foo"}.Username())

	profile := Profile{
		Email: "test@blend.com",
	}

	assert.Equal("test", profile.Username())

	profile = Profile{
		Email: "test2@blendlabs.com",
	}
	assert.Equal("test2", profile.Username())

	profile = Profile{
		Email: "test2+why@blendlabs.com",
	}
	assert.Equal("test2+why", profile.Username())

	profile = Profile{
		Email: "obnoxious@foo@bar@baz.com",
	}
	assert.Equal("obnoxious", profile.Username())
}
