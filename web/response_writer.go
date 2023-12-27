/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"io"
	"net/http"
)

// ResponseWriter is a super-type of http.ResponseWriter that includes
// the StatusCode and ContentLength for the request
type ResponseWriter interface {
	http.Flusher
	http.ResponseWriter
	io.Closer
	StatusCode() int
	ContentLength() int
	InnerResponse() http.ResponseWriter
}
