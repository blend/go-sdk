/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import "github.com/blend/go-sdk/webutil"

// XMLResult is a json result.
type XMLResult struct {
	StatusCode int
	Response   interface{}
}

// Render renders the result
func (ar *XMLResult) Render(ctx *Ctx) error {
	return webutil.WriteXML(ctx.Response, ar.StatusCode, ar.Response)
}
