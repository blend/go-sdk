package r2

import "net/http"

// OptCookie adds a cookie.
func OptCookie(cookie *http.Cookie) Option {
	return func(r *Request) error {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.AddCookie(cookie)
		return nil
	}
}

// OptCookieValue adds a cookie with a given name and value.
func OptCookieValue(name, value string) Option {
	return func(r *Request) error {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.AddCookie(&http.Cookie{Name: name, Value: value})
		return nil
	}
}
