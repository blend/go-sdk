package r2

import (
	"net/http"

	"github.com/blend/go-sdk/ex"
)

// OptKeepAlive enables keep alives.
func OptKeepAlive() Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		if r.Client.Transport == nil {
			return ex.Class("r2; opt keep alive; you must provide a transport to enable keep alives")
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.DisableKeepAlives = false
		} else {
			return ex.Class("r2; opt keep alive; cannot enable keep alives on non http transport")
		}
		return nil
	}
}

// OptDisableKeepAlives disables keep alives.
func OptDisableKeepAlives(disableKeepAlives bool) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.DisableKeepAlives = disableKeepAlives
		}
		return nil
	}
}
