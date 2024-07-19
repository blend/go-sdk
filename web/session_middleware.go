/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

// SessionAware is an action that injects the session into the context.
func SessionAware(action Action) Action {
	return func(ctx *Ctx) Result {
		session, err := ctx.App.Auth.VerifyOrExtendSession(ctx)
		if err != nil && !IsErrSessionInvalid(err) {
			return ctx.DefaultProvider.InternalError(err)
		}
		ctx.Session = session
		ctx.WithContext(WithSession(ctx.Context(), session))
		return action(ctx)
	}
}

// SessionAwareForLogout is an action that injects the session into the context, but does
// not extend it if there is a session lifetime handler on the auth manager.
func SessionAwareForLogout(action Action) Action {
	return func(ctx *Ctx) Result {
		_, session, err := ctx.App.Auth.VerifySession(ctx)
		if err != nil && !IsErrSessionInvalid(err) {
			return ctx.DefaultProvider.InternalError(err)
		}
		ctx.Session = session
		ctx.WithContext(WithSession(ctx.Context(), session))
		return action(ctx)
	}
}

// SessionRequired is an action that requires a session to be present
// or identified in some form on the request.
func SessionRequired(action Action) Action {
	return func(ctx *Ctx) Result {
		session, err := ctx.App.Auth.VerifyOrExtendSession(ctx)
		if err != nil && !IsErrSessionInvalid(err) {
			return ctx.DefaultProvider.InternalError(err)
		}
		if session == nil {
			return ctx.App.Auth.LoginRedirect(ctx)
		}
		ctx.Session = session
		ctx.WithContext(WithSession(ctx.Context(), session))
		return action(ctx)
	}
}

// SessionMiddleware implements a custom notAuthorized action.
func SessionMiddleware(notAuthorized Action) Middleware {
	return func(action Action) Action {
		return func(ctx *Ctx) Result {
			session, err := ctx.App.Auth.VerifyOrExtendSession(ctx)
			if err != nil && !IsErrSessionInvalid(err) {
				return ctx.DefaultProvider.InternalError(err)
			}

			if session == nil {
				if notAuthorized != nil {
					return notAuthorized(ctx)
				}
				return ctx.App.Auth.LoginRedirect(ctx)
			}
			ctx.Session = session
			ctx.WithContext(WithSession(ctx.Context(), session))
			return action(ctx)
		}
	}
}
