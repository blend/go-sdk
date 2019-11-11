package request

import (
	"time"

	"github.com/blend/go-sdk/exception"
)

const (
	// ErrMultipleBodySources is an error returned if a request has both the post body and post data set.
	ErrMultipleBodySources = exception.Class("Cannot set both `Body` and `Post Data`")

	// ErrRequiresTransport is an error michael turner is going to hate.
	ErrRequiresTransport = exception.Class("Request settings require a http.Transport to be provided")
)

const (
	// DefaultKeepAlive returns if we should use a keep alive by default.
	DefaultKeepAlive = false
	// DefaultKeepAliveTimeout is the default time to keep idle connections open before they're closed.
	DefaultKeepAliveTimeout = 60 * time.Second
)

const (
	// MethodGet is a method.
	MethodGet = "GET"
	// MethodPost is a method.
	MethodPost = "POST"
	// MethodPut is a method.
	MethodPut = "PUT"
	// MethodPatch is a method.
	MethodPatch = "PATCH"
	// MethodDelete is a method.
	MethodDelete = "DELETE"
	// MethodOptions is a method.
	MethodOptions = "OPTIONS"
)

const (
	// HeaderConnection is a http header.
	HeaderConnection = "Connection"
	// HeaderContentType is a http header.
	HeaderContentType = "Content-Type"
)

const (
	// ConnectionKeepAlive is a connection header value.
	ConnectionKeepAlive = "keep-alive"
)

const (
	// ContentTypeApplicationJSON is a content type header value.
	ContentTypeApplicationJSON = "application/json; charset=utf-8"
	// ContentTypeApplicationXML is a content type header value.
	ContentTypeApplicationXML = "application/xml"
	// ContentTypeTextXML is a content type header value.
	ContentTypeTextXML = "text/xml"
	// ContentTypeApplicationFormEncoded is a content type header value.
	ContentTypeApplicationFormEncoded = "application/x-www-form-urlencoded"
	// ContentTypeApplicationOctetStream is a content type header value.
	ContentTypeApplicationOctetStream = "application/octet-stream"
)
