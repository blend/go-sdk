package webutil

import (
	"crypto/tls"
	"net/http/httptrace"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestHTTPTraceTrace(t *testing.T) {
	assert := assert.New(t)

	fixedTime := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	ht := &HTTPTrace{NowProvider: func() time.Time {
		return fixedTime
	}}
	ct := ht.Trace()
	assert.NotNil(ct)
	// Invoke `GetConn`.
	assert.True(ht.GetConn.IsZero())
	ct.GetConn("")
	assert.ReferenceEqual(fixedTime, ht.GetConn)
	// Invoke `GotConn`.
	assert.True(ht.GotConn.IsZero())
	ct.GotConn(httptrace.GotConnInfo{})
	assert.ReferenceEqual(fixedTime, ht.GotConn)
	// Invoke `PutIdleConn`.
	assert.True(ht.PutIdleConn.IsZero())
	ct.PutIdleConn(nil)
	assert.ReferenceEqual(fixedTime, ht.PutIdleConn)
	// Invoke `GotFirstResponseByte`.
	d := 573 * time.Millisecond
	ht.WroteRequest = fixedTime.Add(-d)
	assert.True(ht.GotFirstResponseByte.IsZero())
	assert.Equal(0, ht.ServerElapsed)
	ct.GotFirstResponseByte()
	assert.ReferenceEqual(fixedTime, ht.GotFirstResponseByte)
	assert.Equal(d, ht.ServerElapsed)
	ht.WroteRequest = time.Time{} // Reset
	// Invoke `DNSStart`.
	assert.True(ht.DNSStart.IsZero())
	ct.DNSStart(httptrace.DNSStartInfo{})
	assert.ReferenceEqual(fixedTime, ht.DNSStart)
	// Invoke `DNSDone`.
	d = 9302 * time.Millisecond
	ht.DNSStart = fixedTime.Add(-d)
	assert.True(ht.DNSDone.IsZero())
	assert.Equal(0, ht.DNSElapsed)
	ct.DNSDone(httptrace.DNSDoneInfo{})
	assert.ReferenceEqual(fixedTime, ht.DNSDone)
	assert.Equal(d, ht.DNSElapsed)
	ht.DNSStart = time.Time{} // Reset
	// Invoke `ConnectStart`.
	assert.True(ht.ConnectStart.IsZero())
	ct.ConnectStart("", "")
	assert.ReferenceEqual(fixedTime, ht.ConnectStart)
	// Invoke `ConnectDone`.
	d = 1029277 * time.Millisecond
	ht.ConnectStart = fixedTime.Add(-d)
	assert.True(ht.ConnectDone.IsZero())
	assert.Equal(0, ht.DialElapsed)
	ct.ConnectDone("", "", nil)
	assert.ReferenceEqual(fixedTime, ht.ConnectDone)
	assert.Equal(d, ht.DialElapsed)
	ht.ConnectStart = time.Time{} // Reset
	// Invoke `TLSHandshakeStart`.
	assert.True(ht.TLSHandshakeStart.IsZero())
	ct.TLSHandshakeStart()
	assert.ReferenceEqual(fixedTime, ht.TLSHandshakeStart)
	// Invoke `TLSHandshakeDone`.
	d = 8888 * time.Millisecond
	ht.TLSHandshakeStart = fixedTime.Add(-d)
	assert.True(ht.TLSHandshakeDone.IsZero())
	assert.Equal(0, ht.TLSHandshakeElapsed)
	ct.TLSHandshakeDone(tls.ConnectionState{}, nil)
	assert.ReferenceEqual(fixedTime, ht.TLSHandshakeDone)
	assert.Equal(d, ht.TLSHandshakeElapsed)
	ht.TLSHandshakeStart = time.Time{} // Reset
	// Invoke `WroteHeaders`.
	assert.True(ht.WroteHeaders.IsZero())
	ct.WroteHeaders()
	assert.ReferenceEqual(fixedTime, ht.WroteHeaders)
}

func TestHTTPTraceTraceWroteRequest(t *testing.T) {
	assert := assert.New(t)

	fixedTime := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	ht := &HTTPTrace{NowProvider: func() time.Time {
		return fixedTime
	}}
	ct := ht.Trace()
	assert.NotNil(ct)

	// Invoke `WroteRequest`, each of `ConnectDone`, `GetConn`, `GotConn` unset
	assert.True(ht.WroteRequest.IsZero())
	ct.WroteRequest(httptrace.WroteRequestInfo{})
	assert.ReferenceEqual(fixedTime, ht.WroteRequest)
	assert.Equal(0, ht.RequestElapsed)

	// Invoke `WroteRequest`, `ConnectDone` is set
	d := 38308 * time.Millisecond
	ht.ConnectDone = fixedTime.Add(-d)
	ct.WroteRequest(httptrace.WroteRequestInfo{})
	assert.Equal(d, ht.RequestElapsed)
	ht.ConnectDone = time.Time{} // Reset

	// Invoke `WroteRequest`, `GetConn` is set
	d = 1003203 * time.Millisecond
	ht.GetConn = fixedTime.Add(-d)
	ct.WroteRequest(httptrace.WroteRequestInfo{})
	assert.Equal(d, ht.RequestElapsed)
	ht.GetConn = time.Time{} // Reset

	// Invoke `WroteRequest`, `GotConn` is set
	d = 600729 * time.Millisecond
	ht.GotConn = fixedTime.Add(-d)
	ct.WroteRequest(httptrace.WroteRequestInfo{})
	assert.Equal(d, ht.RequestElapsed)
	ht.GotConn = time.Time{} // Reset
}

func TestHTTPTraceNow(t *testing.T) {
	assert := assert.New(t)

	ht := &HTTPTrace{}
	// Using default now()
	before := time.Now().UTC()
	time.Sleep(time.Microsecond)
	now := ht.Now()
	assert.True(now.After(before))
	time.Sleep(time.Microsecond)
	after := time.Now().UTC()
	assert.True(now.Before(after))

	// Using `NowProvider`
	fixedTime := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	ht = &HTTPTrace{NowProvider: func() time.Time {
		return fixedTime
	}}
	now = ht.Now()
	assert.ReferenceEqual(fixedTime, now)
}
