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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNoFollowRedirects(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(http.ErrUseLastResponse, NoFollowRedirects()(nil, nil))

	second := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK!")
	}))
	defer second.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, second.URL, http.StatusTemporaryRedirect)
	}))
	defer server.Close()

	client := http.Client{
		CheckRedirect: NoFollowRedirects(),
	}

	res, err := client.Get(server.URL)
	assert.Nil(err)
	defer res.Body.Close()
	assert.Equal(307, res.StatusCode, "the redirect status code should be returned by the server")
}
