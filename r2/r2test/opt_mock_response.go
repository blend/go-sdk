package r2test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/blend/go-sdk/r2"
)

// OptMockResponseString mocks a string response.
func OptMockResponseString(response string) r2.Option {
	return OptMockResponse(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprint(rw, response)
	}))
}

// OptMockResponseStringStatus mocks a string response with a given status code.
func OptMockResponseStringStatus(statusCode int, response string) r2.Option {
	return OptMockResponse(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(statusCode)
		fmt.Fprint(rw, response)
	}))
}

// OptMockResponse mocks a response by creating an httptest server.
func OptMockResponse(handler http.Handler) r2.Option {
	return func(r *r2.Request) error {
		server := httptest.NewServer(handler)
		parsedURL, _ := url.Parse(server.URL)
		if r.Request.URL == nil {
			// unclear if this is even possible
			r.Request.URL = parsedURL
		} else {
			r.Request.URL.Scheme = parsedURL.Scheme
			r.Request.URL.Host = parsedURL.Host
		}

		if r.Closer != nil {
			originalCloser := r.Closer
			r.Closer = func() error {
				server.Close()
				return originalCloser()
			}
		} else {
			r.Closer = func() error {
				server.Close()
				return nil
			}
		}
		return nil
	}
}
