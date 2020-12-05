package webutil

import (
	"crypto/tls"
	"net/http/httptrace"
	"time"
)

// HTTPTrace is timing information for the full http call.
type HTTPTrace struct {
	GetConn     time.Time `json:"getConn"`
	GotConn     time.Time `json:"gotConn"`
	PutIdleConn time.Time `json:"putIdleConn"`

	DNSStart time.Time `json:"dnsStart"`
	DNSDone  time.Time `json:"dnsDone"`

	ConnectStart time.Time `json:"connectStart"`
	ConnectDone  time.Time `json:"connectDone"`

	TLSHandshakeStart time.Time `json:"tlsHandshakeStart"`
	TLSHandshakeDone  time.Time `json:"tlsHandshakeDone"`

	WroteHeaders         time.Time `json:"wroteHeaders"`
	WroteRequest         time.Time `json:"wroteRequest"`
	GotFirstResponseByte time.Time `json:"gotFirstResponseByte"`

	DNSElapsed          time.Duration `json:"dnsElapsed"`
	TLSHandshakeElapsed time.Duration `json:"tlsHandshakeElapsed"`
	DialElapsed         time.Duration `json:"dialElapsed"`
	RequestElapsed      time.Duration `json:"requestElapsed"`
	ServerElapsed       time.Duration `json:"severElapsed"`

	NowProvider func() time.Time
}

// Trace returns the trace binder.
func (ht *HTTPTrace) Trace() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(_ string) {
			ht.GetConn = ht.Now()
		},
		GotConn: func(_ httptrace.GotConnInfo) {
			ht.GotConn = ht.Now()
		},
		PutIdleConn: func(_ error) {
			ht.PutIdleConn = ht.Now()
		},
		GotFirstResponseByte: func() {
			ht.GotFirstResponseByte = ht.Now()
			ht.ServerElapsed = ht.GotFirstResponseByte.Sub(ht.WroteRequest)
		},
		DNSStart: func(_ httptrace.DNSStartInfo) {
			ht.DNSStart = ht.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			ht.DNSDone = ht.Now()
			ht.DNSElapsed = ht.DNSDone.Sub(ht.DNSStart)
		},
		ConnectStart: func(_, _ string) {
			ht.ConnectStart = ht.Now()
		},
		ConnectDone: func(_, _ string, _ error) {
			ht.ConnectDone = ht.Now()
			ht.DialElapsed = ht.ConnectDone.Sub(ht.ConnectStart)
		},
		TLSHandshakeStart: func() {
			ht.TLSHandshakeStart = ht.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			ht.TLSHandshakeDone = ht.Now()
			ht.TLSHandshakeElapsed = ht.TLSHandshakeDone.Sub(ht.TLSHandshakeStart)
		},
		WroteHeaders: func() {
			ht.WroteHeaders = ht.Now()
		},
		WroteRequest: func(_ httptrace.WroteRequestInfo) {
			ht.WroteRequest = ht.Now()
			if !ht.ConnectDone.IsZero() {
				ht.RequestElapsed = ht.WroteRequest.Sub(ht.ConnectDone)
			} else if !ht.GetConn.IsZero() {
				ht.RequestElapsed = ht.WroteRequest.Sub(ht.GetConn)
			} else if !ht.GotConn.IsZero() {
				ht.RequestElapsed = ht.WroteRequest.Sub(ht.GotConn)
			}
		},
	}
}

// Now returns the current time (UTC); this can be customized with a
// `NowProvider`.
func (ht HTTPTrace) Now() time.Time {
	if ht.NowProvider != nil {
		return ht.NowProvider()
	}
	return time.Now().UTC()
}
