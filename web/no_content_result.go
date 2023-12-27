/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
