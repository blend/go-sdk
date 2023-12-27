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
	"net/url"
	"strings"
)

// MustParseURL parses a url and panics if there is an error.
func MustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}

// URLWithScheme returns a copy url with a given scheme.
func URLWithScheme(u *url.URL, scheme string) *url.URL {
	copy := *u
	copy.Scheme = scheme
	return &copy
}

// URLWithHost returns a copy url with a given host.
func URLWithHost(u *url.URL, host string) *url.URL {
	copy := *u
	copy.Host = host
	return &copy
}

// URLWithPort returns a copy url with a given pprt attached to the host.
func URLWithPort(u *url.URL, port string) *url.URL {
	copy := *u
	host := copy.Host
	if strings.Contains(host, ":") {
		parts := strings.SplitN(host, ":", 2)
		copy.Host = parts[0] + ":" + port
	} else {
		copy.Host = host + ":" + port
	}
	return &copy
}

// URLWithPath returns a copy url with a given path.
func URLWithPath(u *url.URL, path string) *url.URL {
	copy := *u
	copy.Path = path
	return &copy
}

// URLWithRawQuery returns a copy url with a given raw query.
func URLWithRawQuery(u *url.URL, rawQuery string) *url.URL {
	copy := *u
	copy.RawQuery = rawQuery
	return &copy
}

// URLWithQuery returns a copy url with a given raw query.
func URLWithQuery(u *url.URL, key, value string) *url.URL {
	copy := *u
	queryValues := copy.Query()
	queryValues.Add(key, value)
	copy.RawQuery = queryValues.Encode()
	return &copy
}
