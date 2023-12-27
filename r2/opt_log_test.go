/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestOptLog(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	log, err := logger.New(logger.OptOutput(buf), logger.OptAll())
	assert.Nil(err)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, "OK!\n")
	}))
	defer server.Close()

	_, err = New(server.URL, OptLog(log)).Discard()
	assert.Nil(err)
	assert.NotEmpty(buf.String())
}
