/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestDecodeBasicAuth(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", "http://example.invalid", nil)
	assert.Nil(err)

	// No authorization header.
	_, _, err = DecodeBasicAuth(req)
	assert.True(ErrIsUnauthorized(err))

	// Authorization header not in form "Basic ..."
	req.Header.Set(HeaderAuthorization, "not-basic")
	_, _, err = DecodeBasicAuth(req)
	assert.True(ErrIsUnauthorized(err))
	req.Header.Set(HeaderAuthorization, "NotBasic bHVsei1zZWNyZXQ=")
	_, _, err = DecodeBasicAuth(req)
	assert.True(ErrIsUnauthorized(err))

	// With authorization header; invalid base64
	req.Header.Set(HeaderAuthorization, "Basic ???")
	_, _, err = DecodeBasicAuth(req)
	assert.True(ErrIsUnauthorized(err))

	// With authorization header; base64 encoded content not in form `un:pw`
	req.Header.Set(HeaderAuthorization, "Basic bHVsei1zZWNyZXQ=")
	_, _, err = DecodeBasicAuth(req)
	assert.True(ErrIsUnauthorized(err))

	// With authorization header; valid base64 header
	req.Header.Set(HeaderAuthorization, "Basic am9leUBtYWlsLmludmFsaWQ6cHdzMzNrciF0")
	username, password, err := DecodeBasicAuth(req)
	assert.Nil(err)
	assert.Equal(username, "joey@mail.invalid")
	assert.Equal(password, "pws33kr!t")
}
