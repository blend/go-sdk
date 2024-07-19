/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import "net/http"

// HTTPClient is a client that can send http requests.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
