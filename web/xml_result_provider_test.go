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

func TestXMLResultProvider(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(http.StatusNotFound, XML.NotFound().(*XMLResult).StatusCode)
	assert.Equal("Not Found", XML.NotFound().(*XMLResult).Response)

	assert.Equal(http.StatusUnauthorized, XML.NotAuthorized().(*XMLResult).StatusCode)
	assert.Equal("Not Authorized", XML.NotAuthorized().(*XMLResult).Response)

	assert.Equal(http.StatusBadRequest, XML.BadRequest(nil).(*XMLResult).StatusCode)
	assert.Equal("Bad Request", XML.BadRequest(nil).(*XMLResult).Response)

	assert.Equal(http.StatusBadRequest, XML.BadRequest(fmt.Errorf("bad-request")).(*XMLResult).StatusCode)
	assert.Equal(fmt.Errorf("bad-request"), XML.BadRequest(fmt.Errorf("bad-request")).(*XMLResult).Response)

	assert.Equal(http.StatusOK, XML.OK().(*XMLResult).StatusCode)
	assert.Equal("OK!", XML.OK().(*XMLResult).Response)

	assert.Equal(http.StatusBadGateway, XML.Status(http.StatusBadGateway, "test").(*XMLResult).StatusCode)
	assert.Equal("test", XML.Status(http.StatusBadGateway, "test").(*XMLResult).Response)

	assert.Equal(http.StatusOK, XML.Result("foo").(*XMLResult).StatusCode)
	assert.Equal("foo", XML.Result("foo").(*XMLResult).Response)

	internalError := XML.InternalError(fmt.Errorf("only a test"))

	typed, ok := internalError.(*LoggedErrorResult)
	assert.True(ok)
	assert.Equal(fmt.Errorf("only a test"), typed.Error)
	inner := typed.Result.(*XMLResult)
	assert.Equal(http.StatusInternalServerError, inner.StatusCode)
	assert.Equal(fmt.Errorf("only a test"), inner.Response)
}
