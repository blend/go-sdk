/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sanitize

import (
	"net/http"
	"strings"
)

// Request sanitizes a given request.
func Request(r *http.Request, opts ...RequestOption) *http.Request {
	return NewRequestSanitizer(opts...).Sanitize(r)
}

// NewRequestSanitizer creates a new request sanitizer.
func NewRequestSanitizer(opts ...RequestOption) RequestSanitizer {
	r := RequestSanitizer{
		DisallowedHeaders:     DefaultSanitizationDisallowedHeaders,
		DisallowedQueryParams: DefaultSanitizationDisallowedQueryParams,
		KeyValuesSanitizer:    KeyValuesSanitizerFunc(DefaultKeyValuesSanitizerFunc),
		PathSanitizer:         PathSanitizerFunc(DefaultPathSanitizerFunc),
	}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}

// RequestOption is a function that mutates sanitization options.
type RequestOption func(*RequestSanitizer)

// OptRequestAddDisallowedHeaders adds disallowed request headers, augmenting defaults.
func OptRequestAddDisallowedHeaders(headers ...string) RequestOption {
	return func(ro *RequestSanitizer) {
		ro.DisallowedHeaders = append(ro.DisallowedHeaders, headers...)
	}
}

// OptRequestSetDisallowedHeaders sets the disallowed request headers, overwriting defaults.
func OptRequestSetDisallowedHeaders(headers ...string) RequestOption {
	return func(ro *RequestSanitizer) {
		ro.DisallowedHeaders = headers
	}
}

// OptRequestAddDisallowedQueryParams adds disallowed request query params, augmenting defaults.
func OptRequestAddDisallowedQueryParams(queryParams ...string) RequestOption {
	return func(rs *RequestSanitizer) {
		rs.DisallowedQueryParams = append(rs.DisallowedQueryParams, queryParams...)
	}
}

// OptRequestSetDisallowedQueryParams sets the disallowed request query params, overwriting defaults.
func OptRequestSetDisallowedQueryParams(queryParams ...string) RequestOption {
	return func(rs *RequestSanitizer) {
		rs.DisallowedQueryParams = queryParams
	}
}

// OptRequestKeyValuesSanitizer sets the request key values sanitizer.
func OptRequestKeyValuesSanitizer(valueSanitizer KeyValuesSanitizer) RequestOption {
	return func(rs *RequestSanitizer) {
		rs.KeyValuesSanitizer = valueSanitizer
	}
}

// OptRequestPathSanitizer sets the request path sanitizer.
func OptRequestPathSanitizer(pathSanitizer PathSanitizer) RequestOption {
	return func(rs *RequestSanitizer) {
		rs.PathSanitizer = pathSanitizer
	}
}

// RequestSanitizer are options for sanitization of http requests.
type RequestSanitizer struct {
	DisallowedHeaders     []string
	DisallowedQueryParams []string
	KeyValuesSanitizer    KeyValuesSanitizer
	PathSanitizer         PathSanitizer
}

// Sanitize applies sanitization options to a given request.
func (rs RequestSanitizer) Sanitize(r *http.Request) *http.Request {
	if r == nil {
		return nil
	}

	copy := r.Clone(r.Context())
	for header, values := range copy.Header {
		if rs.IsHeaderDisallowed(header) {
			copy.Header[header] = rs.KeyValuesSanitizer.SanitizeKeyValues(header, values...)
		}
	}
	if copy.URL != nil {
		queryParams := copy.URL.Query()
		for queryParam, values := range queryParams {
			if rs.IsQueryParamDisallowed(queryParam) {
				queryParams[queryParam] = rs.KeyValuesSanitizer.SanitizeKeyValues(queryParam, values...)
			}
		}
		copy.URL.RawQuery = queryParams.Encode()

		// also sanitize the path
		copy.URL.Path = rs.PathSanitizer.SanitizePath(copy.URL.Path)
	}
	return copy
}

// IsHeaderDisallowed returns if a header is in the disallowed list.
func (rs RequestSanitizer) IsHeaderDisallowed(header string) bool {
	for _, disallowedHeader := range rs.DisallowedHeaders {
		if strings.EqualFold(disallowedHeader, header) {
			return true
		}
	}
	return false
}

// IsQueryParamDisallowed returns if a query param is in the disallowed list.
func (rs RequestSanitizer) IsQueryParamDisallowed(queryParam string) bool {
	for _, disallowedQueryParam := range rs.DisallowedQueryParams {
		if strings.EqualFold(disallowedQueryParam, queryParam) {
			return true
		}
	}
	return false
}
