/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import "net/http"

// PanicHandler is a handler for panics that also takes an error.
type PanicHandler func(http.ResponseWriter, *http.Request, interface{})
