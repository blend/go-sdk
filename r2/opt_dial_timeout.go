package r2

import (
	"net"
	"net/http"
	"time"
)

// OptDialTimeout sets the dial timeout.
func OptDialTimeout(d time.Duration) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.Dial = (&net.Dialer{
				Timeout: d,
			}).Dial
		}
		return nil
	}
}
