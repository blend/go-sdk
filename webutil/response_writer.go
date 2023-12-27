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

// ResponseWriter is a response writer that also returns the written status.
type ResponseWriter interface {
	http.ResponseWriter
	ContentLength() int
	StatusCode() int
	InnerResponse() http.ResponseWriter
}
