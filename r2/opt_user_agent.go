package r2

import (
	"github.com/blend/go-sdk/webutil"
)

// OptUserAgent sets the user agent header on a request.
// It will initialize the request headers map if it's unset.
// It will overwrite any existing user agent header.
func OptUserAgent(userAgent string) Option {
	return RequestOption(webutil.OptHeaderValue(webutil.HeaderUserAgent, userAgent))
}
