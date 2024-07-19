/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reverseproxy

import (
	"net/http"
)

// TransformRequest modifies an HTTP request. This is intended to be used
// during `Proxy.ServeHTTP()` for custom business logic, e.g. checking if a
// client was included and verified in the request.
type TransformRequest = func(*http.Request)
