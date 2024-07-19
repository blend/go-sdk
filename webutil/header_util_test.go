/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestHeaderAny(t *testing.T) {
	assert := assert.New(t)

	assert.True(HeaderAny(http.Header{"Foo": []string{"bar"}}, "foo", "bar"))
	assert.True(HeaderAny(http.Header{"fuzz": []string{"buzz"}, "Foo": []string{"bar"}}, "foo", "bar"))
	assert.False(HeaderAny(http.Header{"fuzz": []string{"buzz"}, "Foo": []string{"bar"}}, "fuzz", "bar"))
	assert.True(HeaderAny(http.Header{"Foo": []string{"example-string,bar"}}, "foo", "bar"))
	assert.True(HeaderAny(http.Header{"Foo": []string{"bar,example-string"}}, "foo", "bar"))
	assert.True(HeaderAny(http.Header{"fuzz": []string{"buzz"}, "Foo": []string{"bar,example-string"}}, "foo", "bar"))
}
