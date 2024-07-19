/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

// Result is the result of a controller.
type Result interface {
	Render(ctx *Ctx) error
}

// ResultPreRender is a result that has a PreRender step.
type ResultPreRender interface {
	PreRender(ctx *Ctx) error
}

// ResultPostRender is a result that has a PostRender step.
type ResultPostRender interface {
	PostRender(ctx *Ctx) error
}
