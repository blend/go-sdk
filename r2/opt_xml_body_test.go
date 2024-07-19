/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptXMLBody(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptXMLBody(xmlTestCase{Status: "OK!"}))
	assert.NotNil(r.Request.Body)

	contents, err := io.ReadAll(r.Request.Body)
	assert.Nil(err)
	assert.Equal(47, r.Request.ContentLength)
	assert.NotNil(r.Request.GetBody)
	assert.Equal("<xmlTestCase><status>OK!</status></xmlTestCase>", string(contents))
}

func TestOptXMLBodyRedirect(t *testing.T) {
	assert := assert.New(t)

	finalServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.ContentLength != 47 {
			http.Error(rw, fmt.Sprintf("final; invalid content length: %d", r.ContentLength), http.StatusBadRequest)
			return
		}
		var actualBody xmlTestCase
		if err := xml.NewDecoder(r.Body).Decode(&actualBody); err != nil {
			http.Error(rw, "final; invalid xml body", http.StatusBadRequest)
			return
		}
		if actualBody.Status != "OK!" {
			http.Error(rw, "final; invalid status", http.StatusBadRequest)
			return
		}
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, "OK!")
	}))
	defer finalServer.Close()

	// we test the redirect to assert that the .GetBody function works as intended
	redirectServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.ContentLength != 47 {
			http.Error(rw, fmt.Sprintf("redirect; invalid content length: %d", r.ContentLength), http.StatusBadRequest)
			return
		}
		// http.StatusTemporaryRedirect necessary here so the body follows the redirect
		http.Redirect(rw, r, finalServer.URL, http.StatusTemporaryRedirect)
	}))
	defer redirectServer.Close()

	r := New(redirectServer.URL,
		OptXMLBody(xmlTestCase{Status: "OK!"}),
	)
	assert.Equal(47, r.Request.ContentLength)
	assert.NotNil(r.Request.Body)
	assert.NotNil(r.Request.GetBody)

	contents, meta, err := r.Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, string(contents))
}
