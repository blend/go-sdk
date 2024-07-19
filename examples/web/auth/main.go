/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

/*
This example is meant to illustrate the bare minimum required to implement an authenticated web app.
It is meant to be extended considerably, and is not secure as currently formed.
You should investigate specific authentication mechanisms like oauth to do the actual authentication.
Caveat; this will only work if you are deploying a single instance of the app.
*/

func main() {
	app := web.MustNew(
		web.OptLog(logger.All()),
		web.OptAuth(web.NewLocalAuthManager()),
	)

	app.ServeStaticCached("/cached", []string{"_static"}, web.SessionMiddleware(func(ctx *web.Ctx) web.Result {
		return web.Text.NotAuthorized()
	}))
	app.ServeStatic("/static", []string{"_static"}, web.SessionMiddleware(func(ctx *web.Ctx) web.Result {
		return web.Text.NotAuthorized()
	}))
	app.ServeStatic("/static_unauthed", []string{"_static"})

	app.Auth.ValidateHandler = func(_ context.Context, session *web.Session) error {
		if session.UserID == "example-string" {
			return fmt.Errorf("example-string isn't allowed here")
		}
		return nil
	}
	app.Auth.LoginRedirectHandler = web.PathRedirectHandler("/login")

	app.Views.AddLiterals(`{{ define "login" }}<a href="/login/user_valid">Login Valid</a><br/><a href="/login/user_notvalid">Login Invalid</a>{{end}}`)
	app.GET("/login", func(r *web.Ctx) web.Result {
		return r.Views.View("login", nil)
	})

	app.GET("/login/:userID", func(r *web.Ctx) web.Result {
		if r.Session != nil {
			r.App.Log.Debugf("already logged in, redirecting")
			return web.RedirectWithMethodf("GET", "/")
		}

		userID, _ := r.Param("userID")
		if !strings.HasSuffix(userID, "_valid") { //maximum security
			return web.Text.NotAuthorized()
		}
		_, err := r.Auth.Login(userID, r)
		if err != nil {
			return web.Text.InternalError(err)
		}
		return web.Text.Result("Logged In")
	}, web.SessionAware)

	app.GET("/logout", func(r *web.Ctx) web.Result {
		if r.Session == nil {
			return web.Text.Result("Weren't logged in anyway.")
		}
		err := r.Auth.Logout(r)
		if err != nil {
			return web.Text.InternalError(err)
		}
		return web.Text.Result("Logged Out")
	}, web.SessionRequired)

	app.GET("/", func(r *web.Ctx) web.Result {
		return web.Text.Result("Yep")
	}, web.SessionRequired)

	if err := graceful.Shutdown(app); err != nil {
		logger.FatalExit(err)
	}
}
