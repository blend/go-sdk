package secrets

import "github.com/blend/go-sdk/ex"

// Common error codes.
const (
	ErrNotFound     ex.Class = "secrets; not found"
	ErrUnauthorized ex.Class = "secrets; not authorized"
	ErrServerError  ex.Class = "secrets; remote error"
)
