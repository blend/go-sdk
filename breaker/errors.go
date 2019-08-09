package breaker

import "github.com/blend/go-sdk/ex"

var (
	// ErrTooManyRequests is returned when the CB state is half open and the requests count is over the cb maxRequests
	ErrTooManyRequests ex.Class = "too many requests"
	// ErrOpenState is returned when the CB state is open
	ErrOpenState ex.Class = "circuit breaker is open"
)
