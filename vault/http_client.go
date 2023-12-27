/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package vault

import "net/http"

// HTTPClient is a client that can send http requests.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
