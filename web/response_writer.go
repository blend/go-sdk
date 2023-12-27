/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"io"
	"net/http"
)

// ResponseWriter is a super-type of http.ResponseWriter that includes
// the StatusCode and ContentLength for the request
type ResponseWriter interface {
	http.Flusher
	http.ResponseWriter
	io.Closer
	StatusCode() int
	ContentLength() int
	InnerResponse() http.ResponseWriter
}
