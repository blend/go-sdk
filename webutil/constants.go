/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"net/http"
	"regexp"
)

const (
	// TestURL can be used in tests for the URL passed to requests.
	//
	// The URL itself sets `https` as the scheme, `test.invalid` as the host,
	// `/test` as the path, and `query=value` as the querystring.
	//
	// Note: .invalid is a special top level domain that will _never_ be assigned
	// to a real registrant, it is always reserved for testing.
	// See: https://www.iana.org/domains/reserved
	TestURL = "https://test.invalid/test?query=value"
)

// Logger flags
const (
	FlagHTTPRequest = "http.request"
)

// HTTP Method constants (also referred to as 'verbs')
//
// They are aliases for the constants in net/http at this point.
const (
	MethodConnect = http.MethodConnect
	MethodGet     = http.MethodGet
	MethodDelete  = http.MethodDelete
	MethodHead    = http.MethodHead
	MethodPatch   = http.MethodPatch
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodOptions = http.MethodOptions
)

// Header names in canonical form.
var (
	HeaderAccept                  = http.CanonicalHeaderKey("Accept")
	HeaderAcceptEncoding          = http.CanonicalHeaderKey("Accept-Encoding")
	HeaderAllow                   = http.CanonicalHeaderKey("Allow")
	HeaderAuthorization           = http.CanonicalHeaderKey("Authorization")
	HeaderCacheControl            = http.CanonicalHeaderKey("Cache-Control")
	HeaderConnection              = http.CanonicalHeaderKey("Connection")
	HeaderContentEncoding         = http.CanonicalHeaderKey("Content-Encoding")
	HeaderContentLength           = http.CanonicalHeaderKey("Content-Length")
	HeaderContentType             = http.CanonicalHeaderKey("Content-Type")
	HeaderCookie                  = http.CanonicalHeaderKey("Cookie")
	HeaderDate                    = http.CanonicalHeaderKey("Date")
	HeaderETag                    = http.CanonicalHeaderKey("etag")
	HeaderForwarded               = http.CanonicalHeaderKey("Forwarded")
	HeaderServer                  = http.CanonicalHeaderKey("Server")
	HeaderSetCookie               = http.CanonicalHeaderKey("Set-Cookie")
	HeaderStrictTransportSecurity = http.CanonicalHeaderKey("Strict-Transport-Security")
	HeaderUserAgent               = http.CanonicalHeaderKey("User-Agent")
	HeaderVary                    = http.CanonicalHeaderKey("Vary")
	HeaderXContentTypeOptions     = http.CanonicalHeaderKey("X-Content-Type-Options")
	HeaderXForwardedFor           = http.CanonicalHeaderKey("X-Forwarded-For")
	HeaderXForwardedHost          = http.CanonicalHeaderKey("X-Forwarded-Host")
	HeaderXForwardedPort          = http.CanonicalHeaderKey("X-Forwarded-Port")
	HeaderXForwardedProto         = http.CanonicalHeaderKey("X-Forwarded-Proto")
	HeaderXForwardedScheme        = http.CanonicalHeaderKey("X-Forwarded-Scheme")
	HeaderXFrameOptions           = http.CanonicalHeaderKey("X-Frame-Options")
	HeaderXRealIP                 = http.CanonicalHeaderKey("X-Real-IP")
	HeaderXServedBy               = http.CanonicalHeaderKey("X-Served-By")
	HeaderXXSSProtection          = http.CanonicalHeaderKey("X-Xss-Protection")
)

/*
SameSite prevents the browser from sending this cookie along with cross-site requests.
The main goal is mitigate the risk of cross-origin information leakage.
It also provides some protection against cross-site request forgery attacks.
Possible values for the flag are "lax", "strict" or "default".
*/
const (
	SameSiteStrict  = "strict"
	SameSiteLax     = "lax"
	SameSiteDefault = "default"
)

var (
	// Allows for a sub-match of the first value after 'for=' to the next
	// comma, semi-colon or space. The match is case-insensitive.
	// forRegex = regexp.MustCompile(`(?i)(?:for=)([^(;|,| )]+)`)

	// Allows for a sub-match for the first instance of scheme (http|https)
	// prefixed by 'proto='. The match is case-insensitive.
	protoRegex = regexp.MustCompile(`(?i)(?:proto=)(https|http)`)
)

// Well known schemes
const (
	SchemeHTTP  = "http"
	SchemeHTTPS = "https"
	SchemeSPDY  = "spdy"
)

// HSTS Cookie Fields
const (
	HSTSMaxAgeFormat      = "max-age=%d"
	HSTSIncludeSubDomains = "includeSubDomains"
	HSTSPreload           = "preload"
)

// Connection header values.
const (
	// ConnectionKeepAlive is a value for the "Connection" header and
	// indicates the server should keep the tcp connection open
	// after the last byte of the response is sent.
	ConnectionKeepAlive = "keep-alive"
)

const (
	// ContentTypeApplicationJSON is a content type for JSON responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeApplicationJSON = "application/json; charset=utf-8"

	// ContentTypeApplicationXML is a content type header value.
	ContentTypeApplicationXML = "application/xml"

	// ContentTypeApplicationFormEncoded is a content type header value.
	ContentTypeApplicationFormEncoded = "application/x-www-form-urlencoded"

	// ContentTypeApplicationOctetStream is a content type header value.
	ContentTypeApplicationOctetStream = "application/octet-stream"

	// ContentTypeHTML is a content type for html responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeHTML = "text/html; charset=utf-8"

	//ContentTypeXML is a content type for XML responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeXML = "text/xml; charset=utf-8"

	// ContentTypeText is a content type for text responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeText = "text/plain; charset=utf-8"

	// ContentEncodingIdentity is the identity (uncompressed) content encoding.
	ContentEncodingIdentity = "identity"

	// ContentEncodingGZIP is the gzip (compressed) content encoding.
	ContentEncodingGZIP = "gzip"

	// ConnectionClose is the connection value of "close"
	ConnectionClose = "close"
)
