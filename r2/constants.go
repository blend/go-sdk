/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import "github.com/blend/go-sdk/webutil"

const (
	// TestURL can be used in tests for the URL passed to r2.New(...)
	//
	// The URL itself sets `https` as the scheme, `test.invalid` as the host,
	// `/test` as the path, and `query=value` as the querystring.
	//
	// Note: .invalid is a top level special domain that will _never_ be assigned
	// to a real registrant, it is always reserved for testing.
	// See: https://www.iana.org/domains/reserved
	TestURL = webutil.TestURL
)

const (
	// MethodGet is a method.
	MethodGet	= webutil.MethodGet
	// MethodPost is a method.
	MethodPost	= webutil.MethodPost
	// MethodPut is a method.
	MethodPut	= webutil.MethodPut
	// MethodPatch is a method.
	MethodPatch	= webutil.MethodPatch
	// MethodDelete is a method.
	MethodDelete	= webutil.MethodDelete
	// MethodOptions is a method.
	MethodOptions	= webutil.MethodOptions
)

var (
	// HeaderConnection is a http header.
	HeaderConnection	= webutil.HeaderConnection
	// HeaderContentType is a http header.
	HeaderContentType	= webutil.HeaderContentType
)

const (
	// ConnectionKeepAlive is a connection header value.
	ConnectionKeepAlive = webutil.ConnectionKeepAlive
)

const (
	// ContentTypeApplicationJSON is a content type header value.
	ContentTypeApplicationJSON	= webutil.ContentTypeApplicationJSON
	// ContentTypeApplicationXML is a content type header value.
	ContentTypeApplicationXML	= webutil.ContentTypeApplicationXML
	// ContentTypeApplicationFormEncoded is a content type header value.
	ContentTypeApplicationFormEncoded	= webutil.ContentTypeApplicationFormEncoded
	// ContentTypeApplicationOctetStream is a content type header value.
	ContentTypeApplicationOctetStream	= webutil.ContentTypeApplicationOctetStream
)
