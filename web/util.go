package web

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/blend/go-sdk/stringutil"
)

// PathRedirectHandler returns a handler for AuthManager.RedirectHandler based on a path.
func PathRedirectHandler(path string) func(*Ctx) *url.URL {
	return func(ctx *Ctx) *url.URL {
		u := *ctx.Request.URL
		u.Path = path
		return &u
	}
}

// NewSessionID returns a new session id.
// It is not a uuid; session ids are generated using a secure random source.
// SessionIDs are generally 64 bytes.
func NewSessionID() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// NewRequestID returns a pseudo-unique key for a request context.
func NewRequestID() string {
	return stringutil.Random(stringutil.Letters, 12)
}

// Base64URLDecode decodes a base64 string.
func Base64URLDecode(raw string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(raw)
}

// Base64URLEncode base64 encodes data.
func Base64URLEncode(raw []byte) string {
	return base64.URLEncoding.EncodeToString(raw)
}

// ParseInt32 parses an int32.
func ParseInt32(v string) int32 {
	parsed, _ := strconv.Atoi(v)
	return int32(parsed)
}

// NewCookie returns a new name + value pair cookie.
func NewCookie(name, value string) *http.Cookie {
	return &http.Cookie{Name: name, Value: value}
}

// CopyHeaders copies headers.
func CopyHeaders(headers http.Header) http.Header {
	output := make(http.Header)
	for key, values := range headers {
		output[key] = values
	}
	return output
}

// CopySingleHeaders copies headers in single value format.
func CopySingleHeaders(headers map[string]string) http.Header {
	output := make(http.Header)
	for key, value := range headers {
		output[key] = []string{value}
	}
	return output
}

// MergeHeaders merges headers.
func MergeHeaders(headers ...http.Header) http.Header {
	output := make(http.Header)
	for _, header := range headers {
		for key, values := range header {
			output[key] = append(output[key], values...)
		}
	}
	return output
}

// ExtractHost splits a host / port pair (or just a host) and returns the host.
// This is largely borrowed from `net/url.splitHostPort` (as of `go1.13.5`).
func ExtractHost(hostport string) string {
	host := hostport

	colon := strings.LastIndexByte(host, ':')
	if colon != -1 {
		host = host[:colon]
	}

	// If `hostport` is an IPv6 address of the form `[::1]:12801`.
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}

	return host
}
