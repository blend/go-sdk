package webutil

import (
	"net/http"
	"regexp"
)

// Logger flags
const (
	HTTPRequest  = "http.request"
	HTTPResponse = "http.response"
)

// canonical header names.
var (
	// RFC7239 defines a new "Forwarded: " header designed to replace the
	// existing use of X-Forwarded-* headers.
	// e.g. Forwarded: for=192.0.2.60;proto=https;by=203.0.113.43
	HeaderForwarded               = http.CanonicalHeaderKey("Forwarded")
	HeaderXForwardedFor           = http.CanonicalHeaderKey("X-Forwarded-For")
	HeaderXForwardedPort          = http.CanonicalHeaderKey("X-Forwarded-Port")
	HeaderXForwardedHost          = http.CanonicalHeaderKey("X-Forwarded-Host")
	HeaderXForwardedProto         = http.CanonicalHeaderKey("X-Forwarded-Proto")
	HeaderXForwardedScheme        = http.CanonicalHeaderKey("X-Forwarded-Scheme")
	HeaderXRealIP                 = http.CanonicalHeaderKey("X-Real-IP")
	HeaderAcceptEncoding          = http.CanonicalHeaderKey("Accept-Encoding")
	HeaderSetCookie               = http.CanonicalHeaderKey("Set-Cookie")
	HeaderCookie                  = http.CanonicalHeaderKey("Cookie")
	HeaderDate                    = http.CanonicalHeaderKey("Date")
	HeaderCacheControl            = http.CanonicalHeaderKey("Cache-Control")
	HeaderConnection              = http.CanonicalHeaderKey("Connection")
	HeaderContentEncoding         = http.CanonicalHeaderKey("Content-Encoding")
	HeaderContentLength           = http.CanonicalHeaderKey("Content-Length")
	HeaderContentType             = http.CanonicalHeaderKey("Content-Type")
	HeaderUserAgent               = http.CanonicalHeaderKey("User-Agent")
	HeaderServer                  = http.CanonicalHeaderKey("Server")
	HeaderVary                    = http.CanonicalHeaderKey("Vary")
	HeaderXServedBy               = http.CanonicalHeaderKey("X-Served-By")
	HeaderXFrameOptions           = http.CanonicalHeaderKey("X-Frame-Options")
	HeaderXXSSProtection          = http.CanonicalHeaderKey("X-Xss-Protection")
	HeaderXContentTypeOptions     = http.CanonicalHeaderKey("X-Content-Type-Options")
	HeaderStrictTransportSecurity = http.CanonicalHeaderKey("Strict-Transport-Security")
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
	forRegex = regexp.MustCompile(`(?i)(?:for=)([^(;|,| )]+)`)
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

const (
	// ContentTypeApplicationJSON is a content type for JSON responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeApplicationJSON = "application/json; charset=UTF-8"

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

	// ConnectionKeepAlive is a value for the "Connection" header and
	// indicates the server should keep the tcp connection open
	// after the last byte of the response is sent.
	ConnectionKeepAlive = "keep-alive"

	// ContentEncodingIdentity is the identity (uncompressed) content encoding.
	ContentEncodingIdentity = "identity"
	// ContentEncodingGZIP is the gzip (compressed) content encoding.
	ContentEncodingGZIP = "gzip"
)
