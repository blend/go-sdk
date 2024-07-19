/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
	message = ApplyMessageOptions(message, OptMessageChannel("#foo"))
	assert.Equal("#foo", message.Channel)

	assert.Empty(message.IconURL)
	message = ApplyMessageOptions(message, OptMessageIconURL("https://foo.bar.com/icon.png"))
	assert.Equal("https://foo.bar.com/icon.png", message.IconURL)

	assert.Empty(message.IconEmoji)
	message = ApplyMessageOptions(message, OptMessageIconEmoji(":fire:"))
	assert.Equal(":fire:", message.IconEmoji)

	assert.Empty(message.Username)
	message = ApplyMessageOptions(message, OptMessageUsername("example-stringdog"))
	assert.Equal("example-stringdog", message.Username)
}
