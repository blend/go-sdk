/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
