package slack

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMessageOptions(t *testing.T) {
	assert := assert.New(t)

	message := Message{
		Text: "this is only a test",
	}

	assert.Empty(message.Channel)
	message = ApplyMessageOptions(message, WithChannel("#foo"))
	assert.Equal("#foo", message.Channel)

	assert.Empty(message.IconURL)
	message = ApplyMessageOptions(message, WithIconURL("https://foo.bar.com/icon.png"))
	assert.Equal("https://foo.bar.com/icon.png", message.IconURL)

	assert.Empty(message.IconEmoji)
	message = ApplyMessageOptions(message, WithIconEmoji(":fire:"))
	assert.Equal(":fire:", message.IconEmoji)

	assert.Empty(message.Username)
	message = ApplyMessageOptions(message, WithUsername("baileydog"))
	assert.Equal("baileydog", message.Username)
}
