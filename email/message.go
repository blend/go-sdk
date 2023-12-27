/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package email

import (
	"strings"

	"github.com/blend/go-sdk/ex"
)

// Errors
const (
	ErrMessageFieldUnset	ex.Class	= "email; message required field unset"
	ErrMessageFieldNewlines	ex.Class	= "email; message field contains newlines"
)

// MergeMessages merges a given set of messages, such that
// messages that appear last override the options in the output
// message provided those fields are set.
//
// If message A has just the `To` field set, and message B
// has just the `Subject`field set, and message C also has the
// `To` field set, the result of:
//
//	MergeMessages(a,b,c)
//
// Will be `To` of C, and `Subject` of B, with C overwriting
// the `To` of A.
func MergeMessages(messages ...Message) (output Message) {
	for _, message := range messages {
		if message.From != "" {
			output.From = message.From
		}
		if len(message.To) > 0 {
			output.To = message.To
		}
		if len(message.CC) > 0 {
			output.To = message.CC
		}
		if len(message.BCC) > 0 {
			output.To = message.BCC
		}
		if message.Subject != "" {
			output.Subject = message.Subject
		}
		if message.TextBody != "" {
			output.TextBody = message.TextBody
		}
		if message.HTMLBody != "" {
			output.HTMLBody = message.HTMLBody
		}
	}
	return
}

// Message is a message to send via. ses.
type Message struct {
	From		string		`json:"from" yaml:"from"`
	To		[]string	`json:"to" yaml:"to"`
	CC		[]string	`json:"cc" yaml:"cc"`
	BCC		[]string	`json:"bcc" yaml:"bcc"`
	Subject		string		`json:"subject" yaml:"subject"`
	TextBody	string		`json:"textBody" yaml:"textBody"`
	HTMLBody	string		`json:"htmlBody" yaml:"htmlBody"`
}

// IsZero returns if the object is set or not.
func (m Message) IsZero() bool {
	return len(m.To) == 0
}

// Validate checks that a message can be sent.
func (m Message) Validate() error {
	if m.From == "" {
		return ex.New(ErrMessageFieldUnset, ex.OptMessage("field: from"))
	}
	if strings.ContainsAny(m.From, "\r\n") {
		return ex.New(ErrMessageFieldNewlines, ex.OptMessage("field: from"))
	}
	if len(m.To) == 0 {
		return ex.New(ErrMessageFieldUnset, ex.OptMessage("field: to"))
	}
	for index, to := range m.To {
		if strings.ContainsAny(to, "\r\n") {
			return ex.New(ErrMessageFieldNewlines, ex.OptMessagef("field: to[%d]", index))
		}
	}
	for index, cc := range m.CC {
		if strings.ContainsAny(cc, "\r\n") {
			return ex.New(ErrMessageFieldNewlines, ex.OptMessagef("field: cc[%d]", index))
		}
	}
	for index, bcc := range m.BCC {
		if strings.ContainsAny(bcc, "\r\n") {
			return ex.New(ErrMessageFieldNewlines, ex.OptMessagef("field: bcc[%d]", index))
		}
	}
	if strings.ContainsAny(m.Subject, "\r\n") {
		return ex.New(ErrMessageFieldNewlines, ex.OptMessage("field: subject"))
	}
	if len(m.TextBody) == 0 && len(m.HTMLBody) == 0 {
		return ex.New(ErrMessageFieldUnset, ex.OptMessage("fields: textBody and htmlBody"))
	}
	return nil
}
