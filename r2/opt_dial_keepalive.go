package r2

import (
	"net"
	"net/http"
	"time"
)

// OptDialKeepAlive sets the dial keep alive duration.
// Only use this if you know what you're doing, the defaults are typically sufficient.
func OptDialKeepAlive(d time.Duration) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.Dial = (&net.Dialer{
				KeepAlive: d,
			}).Dial
		}
		return nil
	}
}
