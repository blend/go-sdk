package webutil

import "net/http"

// ResponseWriter is a response writer that also returns the written status.
type ResponseWriter interface {
	http.ResponseWriter
	ContentLength() int
	StatusCode() int
	InnerResponse() http.ResponseWriter
}
