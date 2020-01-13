package r2

import (
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// OptHeader sets the request headers.
func OptHeader(headers http.Header) Option {
	return RequestOption(webutil.OptHeader(headers))
}

// OptHeaderValue adds or sets a header value.
func OptHeaderValue(key, value string) Option {
	return RequestOption(webutil.OptHeaderValue(key, value))
}
