/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reverseproxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestRedirectHost(t *testing.T) {
	assert := assert.New(t)

	redirect := HTTPRedirect{
		RedirectScheme:	"spdy",
		RedirectHost:	"redirect-host",
	}
	mockedRedirect := httptest.NewServer(redirect)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	url := fmt.Sprintf("%s/foo", mockedRedirect.URL)
	res, err := client.Get(url)
	assert.Nil(err)
	defer res.Body.Close()

	fullBody, err := io.ReadAll(res.Body)
	assert.Nil(err)

	mockedContents := string(fullBody)
	assert.Equal(http.StatusMovedPermanently, res.StatusCode)

	assert.Contains(mockedContents, "spdy://redirect-host/foo")
}

func TestRedirect(t *testing.T) {
	assert := assert.New(t)

	var redirect HTTPRedirect
	mockedRedirect := httptest.NewServer(redirect)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	urlSuffixes := []string{
		"foo/bar",
		"foo/bar/",
		"foo/bar?test=me",
	}

	for _, urlSuffix := range urlSuffixes {
		url := fmt.Sprintf("%s/%s", mockedRedirect.URL, urlSuffix)
		res, err := client.Get(url)
		assert.Nil(err)
		defer res.Body.Close()

		fullBody, err := io.ReadAll(res.Body)
		assert.Nil(err)

		mockedContents := string(fullBody)
		assert.Equal(http.StatusMovedPermanently, res.StatusCode)

		expectedURL := strings.Replace(url, "http", "https", -1)
		assert.Contains(mockedContents, expectedURL)
	}
}
