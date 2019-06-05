package webutil

import "net/http"

// ResponseWrapper is a type that wraps a response.
type ResponseWrapper interface {
	ContentLength() int
	StatusCode() int
	InnerResponse() http.ResponseWriter
}
