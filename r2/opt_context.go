package r2

import (
	"context"
	"github.com/blend/go-sdk/webutil"
)

// OptContext sets the request context.
func OptContext(ctx context.Context) Option {
	return RequestOption(webutil.OptContext(ctx))
}
