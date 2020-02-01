package r2

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/blend/go-sdk/ex"
)

// OptMaxRedirects tells the http client to only follow a given
// number of redirects, overriding the standard library default of 10.
// Use the companion helper `ErrIsTooManyRedirects` to test if the returned error
// from a call indicates the redirect limit was reached.
func OptMaxRedirects(maxRedirects int) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		r.Client.CheckRedirect = func(r *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return ex.New(http.ErrUseLastResponse)
			}
			return nil
		}
		return nil
	}
}

// ErrIsTooManyRedirects returns if the error is too many redirects.
func ErrIsTooManyRedirects(err error) bool {
	if typed, ok := err.(*url.Error); ok {
		return ex.Is(typed.Err, http.ErrUseLastResponse)
	}
	return false
}

func urlStrings(via []*http.Request) []string {
	var output []string
	for _, req := range via {
		output = append(output, fmt.Sprintf("%s %v", strings.ToUpper(req.Method), req.URL.String()))
	}
	return output
}
