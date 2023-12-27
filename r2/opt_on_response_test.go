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

func TestOptOnResponse(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL,
		OptOnResponse(func(_ *http.Request, _ *http.Response, _ time.Time, _ error) error { return nil }),
		OptOnResponse(func(_ *http.Request, _ *http.Response, _ time.Time, _ error) error { return nil }),
	)
	assert.Len(r.OnResponse, 2)
}
