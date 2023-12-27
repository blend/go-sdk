/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import "github.com/blend/go-sdk/webutil"

// XMLResult is a json result.
type XMLResult struct {
	StatusCode	int
	Response	interface{}
}

// Render renders the result
func (ar *XMLResult) Render(ctx *Ctx) error {
	return webutil.WriteXML(ctx.Response, ar.StatusCode, ar.Response)
}
