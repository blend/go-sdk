/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
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

// Base64URLDecode decodes a base64 string.
func Base64URLDecode(raw string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(raw)
}

// Base64URLEncode base64 encodes data.
func Base64URLEncode(raw []byte) string {
	return base64.URLEncoding.EncodeToString(raw)
}

// NewCookie returns a new name + value pair cookie.
func NewCookie(name, value string) *http.Cookie {
	return &http.Cookie{Name: name, Value: value}
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
			for _, value := range values {
				output.Add(key, value)
			}
		}
	}
	return output
}
