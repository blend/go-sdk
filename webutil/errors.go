package webutil

import "github.com/blend/go-sdk/ex"

// Errors
const (
	ErrInvalidSameSite     ex.Class = "invalid cookie same site string value"
	ErrHijackerUnsupported ex.Class = "ResponseWriter doesn't support Hijacker interface"
)
