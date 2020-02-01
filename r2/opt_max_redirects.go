package r2

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/blend/go-sdk/logger"

	"github.com/blend/go-sdk/ex"
)

// OptMaxRedirects tells the http client to only follow a given
// set of redirects, overriding the standard library default of 10.
// If a maximum number of redirects is reached, an exception of class `http.ErrUseLastResponse`
// will be returned, with the redirect history as the exception message.
// NOTE: that will make the exception message incredibly large.
func OptMaxRedirects(maxRedirects int) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		r.Client.CheckRedirect = func(r *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				r = r.WithContext(logger.WithAnnotations(r.Context(), logger.Annotations{
					"via": urlStrings(via),
				}))
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
