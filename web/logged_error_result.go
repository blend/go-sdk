/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

// ResultWithLoggedError logs an error before it renders the result.
func ResultWithLoggedError(result Result, err error) *LoggedErrorResult {
	return &LoggedErrorResult{
		Error:	err,
		Result:	result,
	}
}

var (
	_	Result			= (*LoggedErrorResult)(nil)
	_	ResultPostRender	= (*LoggedErrorResult)(nil)
)

// LoggedErrorResult is a result that returns an error during the prerender phase.
type LoggedErrorResult struct {
	Result	Result
	Error	error
}

// Render renders the result.
func (ler LoggedErrorResult) Render(ctx *Ctx) error {
	return ler.Result.Render(ctx)
}

// PostRender returns the underlying error.
func (ler LoggedErrorResult) PostRender(ctx *Ctx) error {
	return ler.Error
}
