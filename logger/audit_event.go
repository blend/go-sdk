/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

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
	_	Event		= (*AuditEvent)(nil)
	_	TextWritable	= (*AuditEvent)(nil)
	_	JSONWritable	= (*AuditEvent)(nil)
)

// NewAuditEvent returns a new audit event.
func NewAuditEvent(principal, verb string, options ...AuditEventOption) AuditEvent {
	ae := AuditEvent{
		Principal:	principal,
		Verb:		verb,
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

// NewAuditEventFilter returns a new audit event filter.
func NewAuditEventFilter(filter func(context.Context, AuditEvent) (AuditEvent, bool)) Filter {
	return func(ctx context.Context, e Event) (Event, bool) {
		if typed, isTyped := e.(AuditEvent); isTyped {
			return filter(ctx, typed)
		}
		return e, false
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
	Context		string
	Principal	string
	Verb		string
	Noun		string
	Subject		string
	Property	string
	RemoteAddress	string
	UserAgent	string
	Extra		map[string]string
}

// GetFlag implements Event.
func (e AuditEvent) GetFlag() string	{ return Audit }

// WriteText implements TextWritable.
func (e AuditEvent) WriteText(formatter TextFormatter, wr io.Writer) {
	if len(e.Context) > 0 {
		fmt.Fprint(wr, formatter.Colorize("Context:", ansi.ColorLightBlack))
		fmt.Fprint(wr, e.Context)
		fmt.Fprint(wr, Space)
	}
	if len(e.Principal) > 0 {
		fmt.Fprint(wr, formatter.Colorize("Principal:", ansi.ColorLightBlack))
		fmt.Fprint(wr, e.Principal)
		fmt.Fprint(wr, Space)
	}
	if len(e.Verb) > 0 {
		fmt.Fprint(wr, formatter.Colorize("Verb:", ansi.ColorLightBlack))
		fmt.Fprint(wr, e.Verb)
		fmt.Fprint(wr, Space)
	}
	if len(e.Noun) > 0 {
		fmt.Fprint(wr, formatter.Colorize("Noun:", ansi.ColorLightBlack))
		fmt.Fprint(wr, e.Noun)
		fmt.Fprint(wr, Space)
	}
	if len(e.Subject) > 0 {
		fmt.Fprint(wr, formatter.Colorize("Subject:", ansi.ColorLightBlack))
		fmt.Fprint(wr, e.Subject)
		fmt.Fprint(wr, Space)
	}
	if len(e.Property) > 0 {
		fmt.Fprint(wr, formatter.Colorize("Property:", ansi.ColorLightBlack))
		fmt.Fprint(wr, e.Property)
		fmt.Fprint(wr, Space)
	}
	if len(e.RemoteAddress) > 0 {
		fmt.Fprint(wr, formatter.Colorize("Remote Addr:", ansi.ColorLightBlack))
		fmt.Fprint(wr, e.RemoteAddress)
		fmt.Fprint(wr, Space)
	}
	if len(e.UserAgent) > 0 {
		fmt.Fprint(wr, formatter.Colorize("UA:", ansi.ColorLightBlack))
		fmt.Fprint(wr, e.UserAgent)
		fmt.Fprint(wr, Space)
	}
	if len(e.Extra) > 0 {
		var values []string
		for key, value := range e.Extra {
			values = append(values, fmt.Sprintf("%s%s", formatter.Colorize(key+":", ansi.ColorLightBlack), value))
		}
		fmt.Fprint(wr, strings.Join(values, " "))
	}
}

// Decompose implements Decomposer.
func (e AuditEvent) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"context":	e.Context,
		"principal":	e.Principal,
		"verb":		e.Verb,
		"noun":		e.Noun,
		"subject":	e.Subject,
		"property":	e.Property,
		"remoteAddr":	e.RemoteAddress,
		"ua":		e.UserAgent,
		"extra":	e.Extra,
	}
}
