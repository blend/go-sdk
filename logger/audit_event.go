package logger

import (
	"context"
	"fmt"
	"io"
	"strings"
	"github.com/blend/go-sdk/ansi"
)

// these are compile time assertions
var (
	_ Event        = (*AuditEvent)(nil)
	_ TextWritable = (*AuditEvent)(nil)
	_ JSONWritable = (*AuditEvent)(nil)
)

// NewAuditEvent returns a new audit event.
func NewAuditEvent(principal, verb string, options ...AuditEventOption) AuditEvent {
	ae := AuditEvent{
		Principal: principal,
		Verb:      verb,
	}
	for _, option := range options {
		option(&ae)
	}
	return ae
}

// NewAuditEventListener returns a new audit event listener.
func NewAuditEventListener(listener func(context.Context, AuditEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(AuditEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// AuditEventOption is an option for AuditEvents.
type AuditEventOption func(*AuditEvent)

// OptAuditContext sets a field on an AuditEvent.
func OptAuditContext(value string) AuditEventOption {
	return func(ae *AuditEvent) { ae.Context = value }
}

// OptAuditPrincipal sets a field on an AuditEvent.
func OptAuditPrincipal(value string) AuditEventOption {
	return func(ae *AuditEvent) { ae.Principal = value }
}

// OptAuditVerb sets a field on an AuditEvent.
func OptAuditVerb(value string) AuditEventOption {
	return func(ae *AuditEvent) { ae.Verb = value }
}

// OptAuditNoun sets a field on an AuditEvent.
func OptAuditNoun(value string) AuditEventOption {
	return func(ae *AuditEvent) { ae.Noun = value }
}

// OptAuditSubject sets a field on an AuditEvent.
func OptAuditSubject(value string) AuditEventOption {
	return func(ae *AuditEvent) { ae.Subject = value }
}

// OptAuditProperty sets a field on an AuditEvent.
func OptAuditProperty(value string) AuditEventOption {
	return func(ae *AuditEvent) { ae.Property = value }
}

// OptAuditRemoteAddress sets a field on an AuditEvent.
func OptAuditRemoteAddress(value string) AuditEventOption {
	return func(ae *AuditEvent) { ae.RemoteAddress = value }
}

// OptAuditUserAgent sets a field on an AuditEvent.
func OptAuditUserAgent(value string) AuditEventOption {
	return func(ae *AuditEvent) { ae.UserAgent = value }
}

// OptAuditExtra sets a field on an AuditEvent.
func OptAuditExtra(values map[string]string) AuditEventOption {
	return func(ae *AuditEvent) { ae.Extra = values }
}

// AuditEvent is a common type of event detailing a business action by a subject.
type AuditEvent struct {
	Context       string
	Principal     string
	Verb          string
	Noun          string
	Subject       string
	Property      string
	RemoteAddress string
	UserAgent     string
	Extra         map[string]string
}

// GetFlag implements Event.
func (e AuditEvent) GetFlag() string { return Audit }

// WriteText implements TextWritable.
func (e AuditEvent) WriteText(formatter TextFormatter, wr io.Writer) {
	if len(e.Context) > 0 {
		io.WriteString(wr, formatter.Colorize("Context:", ansi.ColorLightBlack))
		io.WriteString(wr, e.Context)
		io.WriteString(wr, Space)
	}
	if len(e.Principal) > 0 {
		io.WriteString(wr, formatter.Colorize("Principal:", ansi.ColorLightBlack))
		io.WriteString(wr, e.Principal)
		io.WriteString(wr, Space)
	}
	if len(e.Verb) > 0 {
		io.WriteString(wr, formatter.Colorize("Verb:", ansi.ColorLightBlack))
		io.WriteString(wr, e.Verb)
		io.WriteString(wr, Space)
	}
	if len(e.Noun) > 0 {
		io.WriteString(wr, formatter.Colorize("Noun:", ansi.ColorLightBlack))
		io.WriteString(wr, e.Noun)
		io.WriteString(wr, Space)
	}
	if len(e.Subject) > 0 {
		io.WriteString(wr, formatter.Colorize("Subject:", ansi.ColorLightBlack))
		io.WriteString(wr, e.Subject)
		io.WriteString(wr, Space)
	}
	if len(e.Property) > 0 {
		io.WriteString(wr, formatter.Colorize("Property:", ansi.ColorLightBlack))
		io.WriteString(wr, e.Property)
		io.WriteString(wr, Space)
	}
	if len(e.RemoteAddress) > 0 {
		io.WriteString(wr, formatter.Colorize("Remote Addr:", ansi.ColorLightBlack))
		io.WriteString(wr, e.RemoteAddress)
		io.WriteString(wr, Space)
	}
	if len(e.UserAgent) > 0 {
		io.WriteString(wr, formatter.Colorize("UA:", ansi.ColorLightBlack))
		io.WriteString(wr, e.UserAgent)
		io.WriteString(wr, Space)
	}
	if len(e.Extra) > 0 {
		var values []string
		for key, value := range e.Extra {
			values = append(values, fmt.Sprintf("%s%s", formatter.Colorize(key+":", ansi.ColorLightBlack), value))
		}
		io.WriteString(wr, strings.Join(values, " "))
	}
}

// Decompose implements Decomposer.
func (e AuditEvent) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"context":    e.Context,
		"principal":  e.Principal,
		"verb":       e.Verb,
		"noun":       e.Noun,
		"subject":    e.Subject,
		"property":   e.Property,
		"remoteAddr": e.RemoteAddress,
		"ua":         e.UserAgent,
		"extra":      e.Extra,
	}
}
