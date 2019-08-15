package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/blend/go-sdk/webutil"
)

// WithRequestContext sets the request context correctly.
func WithRequestContext(ctx context.Context, req *http.Request) {
	*req = *req.WithContext(ctx)
}

// WithRequestFields sets the request context correctly.
func WithRequestFields(req *http.Request, fields Fields) {
	WithRequestContext(WithFields(req.Context(), fields), req)
}

// HTTPLogged returns a middleware that logs a request.
func HTTPLogged(log Triggerable) webutil.Middleware {
	return func(action http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			start := time.Now()
			w := webutil.NewResponseWriter(rw)
			defer func() {
				responseEvent := NewHTTPResponseEvent(req,
					OptHTTPResponseStatusCode(w.StatusCode()),
					OptHTTPResponseContentLength(w.ContentLength()),
					OptHTTPResponseElapsed(time.Since(start)),
				)
				if w.Header() != nil {
					responseEvent.ContentType = w.Header().Get(webutil.HeaderContentType)
					responseEvent.ContentEncoding = w.Header().Get(webutil.HeaderContentEncoding)
				}
				MaybeTrigger(
					req.Context(),
					log,
					responseEvent,
				)
			}()
			MaybeTrigger(req.Context(), log, NewHTTPRequestEvent(req))
			action(w, req)
		}
	}
}
