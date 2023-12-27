/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import (
	"errors"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestReplacePathParameters(t *testing.T) {
	testCases := []struct {
		name   string
		path   string
		params map[string]string

		expectRes string
		expectErr error
	}{
		{
			name: "happy case",
			path: "/resource/:resource_id/children/:child_id",
			params: map[string]string{
				"resource_id": "123",
				"child_id":    "456",
			},

			expectRes: "/resource/123/children/456",
		},
		{
			name: "does not lead with '/'",
			path: "resource/:resource_id/children/:child_id",
			params: map[string]string{
				"resource_id": "123",
				"child_id":    "456",
			},

			expectRes: "resource/123/children/456",
		},
		{
			name: "params include colon prefix",
			path: "/resource/:resource_id/children/:child_id",
			params: map[string]string{
				":resource_id": "123",
				":child_id":    "456",
			},

			expectRes: "/resource/123/children/456",
		},
		{
			name: "more params than needed",
			path: "/resource/:resource_id/children/:child_id",
			params: map[string]string{
				":resource_id": "123",
				":child_id":    "456",
				":other_id":    "789",
			},

			expectRes: "/resource/123/children/456",
		},
		{
			name: "needed params are missing",
			path: "/resource/:resource_id/children/:child_id",
			params: map[string]string{
				":resource_id": "123",
			},

			expectErr: ErrMissingRouteParameters,
		},
		{
			name: "nil params map",
			path: "/resource/:resource_id/children/:child_id",

			expectErr: ErrMissingRouteParameters,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			its := assert.New(t)
			res, err := ReplacePathParameters(tc.path, tc.params)
			its.True(errors.Is(err, tc.expectErr), tc.expectErr)
			its.Equal(tc.expectRes, res)
		})
	}
}
