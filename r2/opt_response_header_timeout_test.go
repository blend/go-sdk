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
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestOptResponseHeaderTimeout(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptResponseHeaderTimeout(time.Second))
	assert.Equal(time.Second, r.Client.Transport.(*http.Transport).ResponseHeaderTimeout)
}
