/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import "net/http"

// NoFollowRedirects returns an http client redirect delegate that returns the
// http.ErrUseLastResponse error.
// This prevents the net/http Client from following any redirects.
func NoFollowRedirects() func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}
