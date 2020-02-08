package webutil

import "net/http"

// HTTPServerOption is a mutator for an http server.
type HTTPServerOption func(*http.Server) error
