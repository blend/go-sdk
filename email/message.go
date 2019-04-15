package email

import (
	"strings"

	"github.com/blend/go-sdk/ex"
)

// Errors
const (
	ErrMessageFieldUnset    ex.Class = "email; message required field unset"
	ErrMessageFieldNewlines ex.Class = "email; message field contains newlines"
)

// Message is a message to send via. ses.
type Message struct {
	From     string   `json:"from" yaml:"from"`
	To       []string `json:"to" yaml:"to"`
	CC       []string `json:"cc" yaml:"cc"`
	BCC      []string `json:"bcc" yaml:"bcc"`
	Subject  string   `json:"subject" yaml:"subject"`
	TextBody string   `json:"textBody" yaml:"textBody"`
	HTMLBody string   `json:"htmlBody" yaml:"htmlBody"`
}

// Resolve applies extra resolution steps.
func (m *Message) Resolve() error {
	return nil
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
