/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestTextResultProvider(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(http.StatusNotFound, Text.NotFound().(*RawResult).StatusCode)
	assert.Equal("Not Found", string(Text.NotFound().(*RawResult).Response))

	assert.Equal(http.StatusUnauthorized, Text.NotAuthorized().(*RawResult).StatusCode)
	assert.Equal("Not Authorized", string(Text.NotAuthorized().(*RawResult).Response))

	assert.Equal(http.StatusBadRequest, Text.BadRequest(nil).(*RawResult).StatusCode)
	assert.Equal("Bad Request", string(Text.BadRequest(nil).(*RawResult).Response))

	assert.Equal(http.StatusBadRequest, Text.BadRequest(fmt.Errorf("bad-request")).(*RawResult).StatusCode)
	assert.Equal("Bad Request: bad-request", string(Text.BadRequest(fmt.Errorf("bad-request")).(*RawResult).Response))

	assert.Equal(http.StatusOK, Text.OK().(*RawResult).StatusCode)
	assert.Equal("OK!", string(Text.OK().(*RawResult).Response))

	assert.Equal(http.StatusBadGateway, Text.Status(http.StatusBadGateway, "test").(*RawResult).StatusCode)
	assert.Equal("test", string(Text.Status(http.StatusBadGateway, "test").(*RawResult).Response))

	assert.Equal(http.StatusOK, Text.Result("foo").(*RawResult).StatusCode)
	assert.Equal("foo", string(Text.Result("foo").(*RawResult).Response))

	internalError := Text.InternalError(fmt.Errorf("only a test"))

	typed, ok := internalError.(*LoggedErrorResult)
	assert.True(ok)
	assert.Equal(fmt.Errorf("only a test"), typed.Error)
	inner := typed.Result.(*RawResult)
	assert.Equal(http.StatusInternalServerError, inner.StatusCode)
	assert.Equal("only a test", string(inner.Response))
}
