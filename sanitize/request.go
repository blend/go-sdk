package sanitize

import (
	"net/http"
	"strings"
)

// Request applies sanitization options to a given request.
func Request(r *http.Request, opts ...RequestOption) *http.Request {
	if r == nil {
		return nil
	}

	options := RequestOptions{
		DisallowedHeaders:     DefaultSanitizationDisallowedHeaders,
		DisallowedQueryParams: DefaultSanitizationDisallowedQueryParams,
		ValueSanitizer:        DefaultValueSanitizer,
	}
	for _, opt := range opts {
		opt(&options)
	}

	copy := r.Clone(r.Context())
	for header, values := range copy.Header {
		if options.IsHeaderDisallowed(header) {
			copy.Header[header] = options.ValueSanitizer(header, values...)
		}
	}

	if copy.URL != nil {
		queryParams := copy.URL.Query()
		for queryParam, values := range queryParams {
			if options.IsQueryParamDisallowed(queryParam) {
				queryParams[queryParam] = options.ValueSanitizer(queryParam, values...)
			}
		}
		copy.URL.RawQuery = queryParams.Encode()
	}

	return copy
}

// OptRequestAddDisallowedHeaders adds disallowed headers, augmenting defaults.
func OptRequestAddDisallowedHeaders(headers ...string) RequestOption {
	return func(ro *RequestOptions) {
		ro.DisallowedHeaders = append(ro.DisallowedHeaders, headers...)
	}
}

// OptRequestSetDisallowedHeaders sets the disallowed headers, overwriting defaults.
func OptRequestSetDisallowedHeaders(headers ...string) RequestOption {
	return func(ro *RequestOptions) {
		ro.DisallowedHeaders = headers
	}
}

// OptRequestAddDisallowedQueryParams adds disallowed query params, augmenting defaults.
func OptRequestAddDisallowedQueryParams(queryParams ...string) RequestOption {
	return func(ro *RequestOptions) {
		ro.DisallowedQueryParams = append(ro.DisallowedQueryParams, queryParams...)
	}
}

// OptRequestSetDisallowedQueryParams sets the disallowed query params, overwriting defaults.
func OptRequestSetDisallowedQueryParams(queryParams ...string) RequestOption {
	return func(ro *RequestOptions) {
		ro.DisallowedQueryParams = queryParams
	}
}

// OptRequestValueSanitizer sets the value sanitizer.
func OptRequestValueSanitizer(valueSanitizer ValueSanitizer) RequestOption {
	return func(ro *RequestOptions) {
		ro.ValueSanitizer = valueSanitizer
	}
}

// RequestOptions are options for sanitization of http requests.
type RequestOptions struct {
	DisallowedHeaders     []string
	DisallowedQueryParams []string
	ValueSanitizer        ValueSanitizer
}

// IsHeaderDisallowed returns if a header is in the disallowed list.
func (ro RequestOptions) IsHeaderDisallowed(header string) bool {
	for _, disallowedHeader := range ro.DisallowedHeaders {
		if strings.EqualFold(disallowedHeader, header) {
			return true
		}
	}
	return false
}

// IsQueryParamDisallowed returns if a query param is in the disallowed list.
func (ro RequestOptions) IsQueryParamDisallowed(queryParam string) bool {
	for _, disallowedQueryParam := range ro.DisallowedQueryParams {
		if strings.EqualFold(disallowedQueryParam, queryParam) {
			return true
		}
	}
	return false
}

// RequestOption is a function that mutates sanitization options.
type RequestOption func(*RequestOptions)
