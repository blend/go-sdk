package web

import (
	"context"
	"net/http"
	"time"
)

// WithTimeout injects the context for a given action with a timeout context.
func WithTimeout(d time.Duration) Middleware {
	return func(action Action) Action {
		return func(r *Ctx) Result {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()

			r.Request = r.Request.WithContext(ctx)

			panicChan := make(chan interface{}, 1)
			resultChan := make(chan Result, 1)

			go func() {
				defer func() {
					if p := recover(); p != nil {
						panicChan <- p
					}
				}()
				resultChan <- action(r)
			}()

			select {
			case p := <-panicChan:
				panic(p)
			case res := <-resultChan:
				return res
			case <-ctx.Done():
				return r.DefaultProvider.Status(http.StatusServiceUnavailable, nil)
			}
		}
	}
}
