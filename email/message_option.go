/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package email

// ApplyMessageOptions applies options to a message and returns
// the mutated copy.
func ApplyMessageOptions(m Message, options ...MessageOption) Message {
	for _, option := range options {
		option(&m)
	}
	return m
}

// MessageOption is a mutator for messages.
type MessageOption func(m *Message)

// OptFrom sets the from address for a message.
func OptFrom(from string) MessageOption {
	return func(m *Message) {
		m.From = from
	}
}

// OptTo sets the to address for a message.
func OptTo(to ...string) MessageOption {
	return func(m *Message) {
		m.To = to
	}
}

// OptCC sets the cc addresses for a message.
func OptCC(cc ...string) MessageOption {
	return func(m *Message) {
		m.CC = cc
	}
}

// OptBCC sets the bcc addresses for a message.
func OptBCC(bcc ...string) MessageOption {
	return func(m *Message) {
		m.BCC = bcc
	}
}

// OptSubject sets the subject for a message.
func OptSubject(subject string) MessageOption {
	return func(m *Message) {
		m.Subject = subject
	}
}

// OptTextBody sets the text body for a message.
func OptTextBody(textBody string) MessageOption {
	return func(m *Message) {
		m.TextBody = textBody
	}
}

// OptHTMLBody sets the html body for a message.
func OptHTMLBody(htmlBody string) MessageOption {
	return func(m *Message) {
		m.HTMLBody = htmlBody
	}
}
