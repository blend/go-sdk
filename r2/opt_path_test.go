/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptPath(t *testing.T) {
	its := assert.New(t)

	r := New(TestURL, OptPath("/not-foo"))
	its.Nil(r.Err)
	its.Equal("/not-foo", r.Request.URL.Path)

	var unset Request
	its.NotNil(OptPath("/not-foo")(&unset))
}

func TestOptPathf(t *testing.T) {
	its := assert.New(t)

	r := New(TestURL, OptPathf("/not-foo/%s", "bar"))
	its.Nil(r.Err)
	its.Equal("/not-foo/bar", r.Request.URL.Path)

	var unset Request
	its.NotNil(OptPathf("/not-foo/%s", "bar")(&unset))
}

func TestOptPathParameterized(t *testing.T) {
	its := assert.New(t)

	r := New(TestURL, OptPathParameterized("resource/:resource_id", map[string]string{"resource_id": "1234"}))
	its.Nil(r.Err)
	its.Equal("/resource/1234", r.Request.URL.Path)
	its.Equal("/resource/:resource_id", GetParameterizedPath(r.Request.Context()))

	var unset Request
	its.NotNil(OptPathParameterized("resource/:resource_id", map[string]string{"resource_id": "1234"})(&unset))

	its.NotNil(OptPathParameterized("resource/:resource_id", map[string]string{})(r))
	its.Nil(OptPathParameterized("resource/:resource_id", map[string]string{"resource_id": "1234", "other_id": "5678"})(r))
}
