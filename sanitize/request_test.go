/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sanitize

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSanitizeRequest(t *testing.T) {
	it := assert.New(t)

	req := &http.Request{
		Header: http.Header{
			"Accept":         {"application/json"},
			"Authorization":  {"Bearer foo", "Bearer bar"},
			"X-Secret-Token": {"super_secret_token"},
		},
		URL: &url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/api/sensitive",
			RawQuery: (url.Values{
				"ok":           {"ok0", "ok1"},
				"access_token": {"super_secret"},
				"sensitive":    {"sensitive0", "sensitive1"},
			}).Encode(),
		},
	}

	sanitizer := NewRequestSanitizer(
		OptRequestAddDisallowedHeaders("X-Secret-Token"),
		OptRequestAddDisallowedQueryParams("sensitive"),
		OptRequestKeyValuesSanitizer(KeyValuesSanitizerFunc(func(key string, values ...string) []string {
			return []string{"***"}
		})),
	)
	output := sanitizer.Sanitize(req)

	it.NotNil(output)
	it.Equal([]string{"application/json"}, req.Header["Accept"])
	it.Equal([]string{"***"}, output.Header["Authorization"])
	it.Equal([]string{"Bearer foo", "Bearer bar"}, req.Header["Authorization"])
	it.Equal([]string{"***"}, output.Header["X-Secret-Token"])

	it.Equal([]string{"***"}, output.URL.Query()["access_token"])
	it.Equal([]string{"***"}, output.URL.Query()["sensitive"])
}
