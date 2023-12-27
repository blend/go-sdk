/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package slack

// NewMessage creates a new message with a given set of options.
func NewMessage(options ...MessageOption) *Message {
	var m Message
	for _, option := range options {
		option(&m)
	}
	return &m
}

// Message is a message sent to slack.
type Message struct {
	Username        string              `json:"username,omitempty"`
	Channel         string              `json:"channel,omitempty"`
	Parse           string              `json:"parse,omitempty"`
	ResponseType    string              `json:"response_type,omitempty"`
	Text            string              `json:"text"`
	IconEmoji       string              `json:"icon_emoji,omitempty"`
	IconURL         string              `json:"icon_url,omitempty"`
	ThreadTimestamp string              `json:"thread_ts,omitempty"`
	AsUser          bool                `json:"as_user,omitempty"`
	LinkNames       bool                `json:"link_names"`
	Attachments     []MessageAttachment `json:"attachments"`

	// Response-specific fields
	BotID     string `json:"bot_id,omitempty"`
	Type      string `json:"type,omitempty"`
	SubType   string `json:"subtype,omitempty"`
	Timestamp string `json:"ts,omitempty"`
}

// MessageAttachment is an attachment for a message.
type MessageAttachment struct {
	Title      string                   `json:"title,omitempty"`
	Color      string                   `json:"color,omitempty"`
	Pretext    string                   `json:"pretext,omitempty"`
	Text       string                   `json:"text,omitempty"`
	MarkdownIn []string                 `json:"mrkdwn_in,omitempty"`
	Fields     []MessageAttachmentField `json:"fields,omitempty"`
}

// MessageAttachmentField is a field on an attachment.
type MessageAttachmentField struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short"`
}
