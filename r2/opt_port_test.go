/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptPort(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptPort(8443))

	expected, _ := url.Parse(TestURL)
	expected.Host = fmt.Sprintf("%s:%s", expected.Hostname(), "8443")
	assert.Equal(expected.String(), r.Request.URL.String())
}
