/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCtxOption(t *testing.T) {
	assert := assert.New(t)

	var ctx Ctx
	assert.Nil(ctx.App)
	OptCtxApp(&App{})(&ctx)
	assert.NotNil(ctx.App)

	assert.Empty(ctx.Auth.CookieDefaults.Name)
	OptCtxAuth(AuthManager{CookieDefaults: http.Cookie{Name: "foo"}})(&ctx)
	assert.Equal("foo", ctx.Auth.CookieDefaults.Name)

	assert.Nil(ctx.DefaultProvider)
	OptCtxDefaultProvider(JSON)(&ctx)
	assert.NotNil(ctx.DefaultProvider)

	assert.Nil(ctx.Views)
	OptCtxViews(&ViewCache{})(&ctx)
	assert.NotNil(ctx.Views)

	assert.Nil(ctx.State)
	OptCtxState(&SyncState{})(&ctx)
	assert.NotNil(ctx.State)

	assert.Nil(ctx.Session)
	OptCtxSession(&Session{})(&ctx)
	assert.NotNil(ctx.Session)

	assert.Nil(ctx.Route)
	OptCtxRoute(&Route{})(&ctx)
	assert.NotNil(ctx.Route)

	assert.Nil(ctx.RouteParams)
	OptCtxRouteParams(RouteParameters{})(&ctx)
	assert.NotNil(ctx.RouteParams)

	assert.Nil(ctx.Tracer)
	OptCtxTracer(&mockTracer{})(&ctx)
	assert.NotNil(ctx.Tracer)
}
