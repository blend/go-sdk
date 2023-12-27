/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestGetRawURLParameterized(t *testing.T) {
	testCases := []struct {
		name	string
		req	*http.Request
		expect	string
	}{
		{
			name:	"nil request",
			req:	nil,
			expect:	"",
		},
		{
			name:	"nil request URL",
			req:	&http.Request{},
			expect:	"",
		},
		{
			name:	"r2 request without any options",
			req:	New("https://example.test/resource/1234").Request,
			expect:	"https://example.test/resource/1234",
		},
		{
			name:	"r2 request using OptPath",
			req:	New("https://example.test", OptPath("resource/1234")).Request,
			expect:	"https://example.test/resource/1234",
		},
		{
			name:	"r2 request using OptPathf",
			req:	New("https://example.test", OptPathf("resource/%d", 1234)).Request,
			expect:	"https://example.test/resource/1234",
		},
		{
			name:	"r2 request using OptPathParameterized",
			req: New("https://example.test", OptPathParameterized("resource/:resource_id/:child_id", map[string]string{
				"resource_id":	"1234",
				"child_id":	"5678",
			})).Request,
			expect:	"https://example.test/resource/:resource_id/:child_id",
		},
		{
			name:	"request with context containing parameterizedPath{}",
			req: (&http.Request{
				URL: &url.URL{},
			}).WithContext(
				WithParameterizedPath(context.Background(), "resource/:resource_id/:child_id"),
			),
			expect:	"/resource/:resource_id/:child_id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			its := assert.New(t)
			its.Equal(tc.expect, GetRawURLParameterized(tc.req))
		})
	}
}
