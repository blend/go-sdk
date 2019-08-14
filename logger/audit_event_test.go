package logger

import (
	"bytes"
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestAuditEventMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	ae := NewAuditEvent(
		"bailey",
		"pooped",
	)

	contents := ae.Decompose()
	assert.NotEmpty(contents)
	assert.Equal("bailey", contents["principal"])
	assert.Equal("pooped", contents["verb"])
}

func TestAuditEventOptions(t *testing.T) {
	assert := assert.New(t)

	ae := NewAuditEvent(
		"bailey",
		"pooped",
		OptAuditContext("event context"),
		OptAuditPrincipal("not bailey"),
		OptAuditVerb("not pooped"),
		OptAuditNoun("audit noun"),
		OptAuditSubject("audit subject"),
		OptAuditProperty("audit property"),
		OptAuditRemoteAddress("remote address"),
		OptAuditUserAgent("user agent"),
		OptAuditExtra(map[string]string{"foo": "bar"}),
	)

	assert.Equal("event context", ae.Context)
	assert.Equal("not bailey", ae.Principal)
	assert.Equal("not pooped", ae.Verb)
	assert.Equal("audit noun", ae.Noun)
	assert.Equal("audit subject", ae.Subject)
	assert.Equal("audit property", ae.Property)
	assert.Equal("remote address", ae.RemoteAddress)
	assert.Equal("user agent", ae.UserAgent)
	assert.Equal("bar", ae.Extra["foo"])
}

func TestAuditEventWriteText(t *testing.T) {
	assert := assert.New(t)

	ae := NewAuditEvent(
		"bailey",
		"pooped",
		OptAuditContext("event context"),
		OptAuditPrincipal("not bailey"),
		OptAuditVerb("not pooped"),
		OptAuditNoun("audit noun"),
		OptAuditSubject("audit subject"),
		OptAuditProperty("audit property"),
		OptAuditRemoteAddress("remote address"),
		OptAuditUserAgent("user agent"),
		OptAuditExtra(map[string]string{"foo": "bar"}),
	)

	buf := new(bytes.Buffer)
	noColor := TextOutputFormatter{
		NoColor: true,
	}

	ae.WriteText(noColor, buf)

	assert.Equal("Context:event context Principal:not bailey Verb:not pooped Noun:audit noun Subject:audit subject Property:audit property Remote Addr:remote address UA:user agent foo:bar", buf.String())
}

func TestAuditEventListener(t *testing.T) {
	assert := assert.New(t)

	ae := NewAuditEvent("bailey", "pooped")

	var didCall bool
	ml := NewAuditEventListener(func(ctx context.Context, ae AuditEvent) {
		didCall = true
	})

	ml(context.Background(), ae)
	assert.True(didCall)
}
