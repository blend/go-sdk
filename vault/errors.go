package vault

import "github.com/blend/go-sdk/ex"

// Common error codes.
const (
	ErrNotFound     ex.Class = "vault; not found"
	ErrUnauthorized ex.Class = "vault; not authorized"
	ErrServerError  ex.Class = "vault; remote error"
)
