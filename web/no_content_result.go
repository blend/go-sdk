/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import "net/http"

var (
	// NoContent is a static result.
	NoContent NoContentResult
)

// NoContentResult returns a no content response.
type NoContentResult struct{}

// Render renders a static result.
func (ncr NoContentResult) Render(ctx *Ctx) error {
	ctx.Response.WriteHeader(http.StatusNoContent)
	return nil
}
