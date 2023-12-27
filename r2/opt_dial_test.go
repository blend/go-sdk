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

func TestOptDial(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptDial(OptDialTimeout(time.Second)))
	assert.NotNil(r.Client.Transport.(*http.Transport).DialContext)
}
