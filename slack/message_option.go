/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package slack

// ApplyMessageOptions applies a set of options against a message and returns the mutated copy.
func ApplyMessageOptions(m Message, options ...MessageOption) Message {
	for _, option := range options {
		option(&m)
	}
	return m
}

// MessageOption is a mutator for messages.
type MessageOption func(m *Message)

// OptMessageChannel sets the channel.
func OptMessageChannel(channel string) MessageOption {
	return func(m *Message) {
		m.Channel = channel
	}
}

// OptMessageChannelOrDefault sets the channel if its unset.
func OptMessageChannelOrDefault(channel string) MessageOption {
	return func(m *Message) {
		if len(m.Channel) == 0 {
			m.Channel = channel
		}
	}
}

// OptMessageResponseType sets the response type.
func OptMessageResponseType(responseType string) MessageOption {
	return func(m *Message) {
		m.ResponseType = responseType
	}
}

// OptMessageUsername sets the username.
func OptMessageUsername(username string) MessageOption {
	return func(m *Message) {
		m.Username = username
	}
}

// OptMessageUsernameOrDefault sets the username.
func OptMessageUsernameOrDefault(username string) MessageOption {
	return func(m *Message) {
		if len(m.Username) == 0 {
			m.Username = username
		}
	}
}

// OptMessageIconEmoji sets the icon emoji.
func OptMessageIconEmoji(emoji string) MessageOption {
	return func(m *Message) {
		m.IconEmoji = emoji
	}
}

// OptMessageIconEmojiOrDefault sets the icon emoji.
func OptMessageIconEmojiOrDefault(emoji string) MessageOption {
	return func(m *Message) {
		if len(m.IconEmoji) == 0 {
			m.IconEmoji = emoji
		}
	}
}

// OptMessageIconURL sets the icon url.
func OptMessageIconURL(url string) MessageOption {
	return func(m *Message) {
		m.IconURL = url
	}
}

// OptMessageIconURLOrDefault sets the icon url.
func OptMessageIconURLOrDefault(url string) MessageOption {
	return func(m *Message) {
		if len(m.IconURL) == 0 {
			m.IconURL = url
		}
	}
}

// OptMessageAttachment adds a message attachment.
func OptMessageAttachment(attachment MessageAttachment) MessageOption {
	return func(m *Message) {
		m.Attachments = append(m.Attachments, attachment)
	}
}
