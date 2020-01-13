package r2

import (
	"net/url"
	"github.com/blend/go-sdk/webutil"
)

// OptQuery sets the full querystring.
func OptQuery(query url.Values) Option {
	return RequestOption(webutil.OptQuery(query))
}

// OptQueryValue adds or sets a query value.
func OptQueryValue(key, value string) Option {
	return RequestOption(webutil.OptQueryValue(key, value))
}
