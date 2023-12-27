/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package reverseproxy

import (
	"net/http"
)

// TransformRequest modifies an HTTP request. This is intended to be used
// during `Proxy.ServeHTTP()` for custom business logic, e.g. checking if a
// client was included and verified in the request.
type TransformRequest = func(*http.Request)
