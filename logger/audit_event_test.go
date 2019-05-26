package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestAuditEventMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	ae := NewAuditEvent(
		"bailey",
		"pooped",
		OptAuditEventMetaOptions(OptEventMetaTimestamp(time.Date(2016, 01, 02, 03, 04, 05, 06, time.UTC))),
	)

	contents, err := json.Marshal(ae)
	assert.Nil(err)

	assert.Contains(string(contents), "bailey")
	assert.Contains(string(contents), "pooped")

	assert.True(strings.HasPrefix(string(contents), `{"_timestamp":"2016-01-02T03:04:05`), string(contents))
}

func TestAuditEventOptions(t *testing.T) {
	assert := assert.New(t)

	ae := NewAuditEvent(
		"bailey",
		"pooped",
		OptAuditEventMetaOptions(OptEventMetaTimestamp(time.Date(2016, 01, 02, 03, 04, 05, 06, time.UTC))),
		OptAuditEventContext("event context"),
		OptAuditEventPrincipal("not bailey"),
		OptAuditEventVerb("not pooped"),
		OptAuditEventNoun("audit noun"),
		OptAuditEventSubject("audit subject"),
		OptAuditEventProperty("audit property"),
		OptAuditEventRemoteAddress("remote address"),
		OptAuditEventUserAgent("user agent"),
		OptAuditEventExtra(map[string]string{"foo": "bar"}),
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
		OptAuditEventMetaOptions(OptEventMetaTimestamp(time.Date(2016, 01, 02, 03, 04, 05, 06, time.UTC))),
		OptAuditEventContext("event context"),
		OptAuditEventPrincipal("not bailey"),
		OptAuditEventVerb("not pooped"),
		OptAuditEventNoun("audit noun"),
		OptAuditEventSubject("audit subject"),
		OptAuditEventProperty("audit property"),
		OptAuditEventRemoteAddress("remote address"),
		OptAuditEventUserAgent("user agent"),
		OptAuditEventExtra(map[string]string{"foo": "bar"}),
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
	ml := NewAuditEventListener(func(ctx context.Context, ae *AuditEvent) {
		didCall = true
	})

	ml(context.Background(), ae)
	assert.True(didCall)
}
