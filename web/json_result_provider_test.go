/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestJSONResultProvider(t *testing.T) {
	assert := assert.New(t)

	notFound, ok := JSON.NotFound().(*JSONResult)
	assert.True(ok)
	assert.Equal(http.StatusNotFound, notFound.StatusCode)
	assert.Equal("Not Found", notFound.Response)

	notAuthorized, ok := JSON.NotAuthorized().(*JSONResult)
	assert.True(ok)
	assert.Equal(http.StatusUnauthorized, notAuthorized.StatusCode)
	assert.Equal("Not Authorized", notAuthorized.Response)

	forbidden, ok := JSON.Forbidden().(*JSONResult)
	assert.True(ok)
	assert.Equal(http.StatusForbidden, forbidden.StatusCode)
	assert.Equal("Forbidden", forbidden.Response)

	badRequest, ok := JSON.BadRequest(nil).(*JSONResult)
	assert.True(ok)
	assert.Equal(http.StatusBadRequest, badRequest.StatusCode)
	assert.Equal("Bad Request", badRequest.Response)

	badRequestErr, ok := JSON.BadRequest(fmt.Errorf("bad-request")).(*JSONResult)
	assert.True(ok)
	assert.Equal(http.StatusBadRequest, badRequestErr.StatusCode)
	assert.Equal("bad-request", badRequestErr.Response)

	okRes, ok := JSON.OK().(*JSONResult)
	assert.True(ok)
	assert.Equal(http.StatusOK, okRes.StatusCode)
	assert.Equal("OK!", okRes.Response)

	statusRes, ok := JSON.Status(http.StatusBadGateway, "test").(*JSONResult)
	assert.True(ok)
	assert.Equal(http.StatusBadGateway, statusRes.StatusCode)
	assert.Equal("test", statusRes.Response)

	res, ok := JSON.Result("foo").(*JSONResult)
	assert.True(ok)
	assert.Equal(http.StatusOK, res.StatusCode)
	assert.Equal("foo", res.Response)

	internalError := JSON.InternalError(fmt.Errorf("only a test"))

	typed, ok := internalError.(*LoggedErrorResult)
	assert.True(ok)
	assert.Equal(fmt.Errorf("only a test"), typed.Error)
	inner := typed.Result.(*JSONResult)
	assert.Equal(http.StatusInternalServerError, inner.StatusCode)
	assert.Equal("only a test", inner.Response)
}
