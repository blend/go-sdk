package web

// SessionAware is an action that injects the session into the context, it acquires a read lock on session.
func SessionAware(action Action) Action {
	return func(ctx *Ctx) Result {
		session, err := ctx.App.Auth.VerifySession(ctx)
		if err != nil && !IsErrSessionInvalid(err) {
			return ctx.DefaultProvider.InternalError(err)
		}
		ctx.Session = session
		return action(ctx)
	}
}

// SessionRequired is an action that requires a session to be present
// or identified in some form on the request, and acquires a read lock on session.
func SessionRequired(action Action) Action {
	return func(ctx *Ctx) Result {
		session, err := ctx.App.Auth.VerifySession(ctx)
		if err != nil && !IsErrSessionInvalid(err) {
			return ctx.DefaultProvider.InternalError(err)
		}
		if session == nil {
			return ctx.App.Auth.LoginRedirect(ctx)
		}
		ctx.Session = session
		return action(ctx)
	}
}

// SessionMiddleware implements a custom notAuthorized action.
func SessionMiddleware(notAuthorized Action) Middleware {
	return func(action Action) Action {
		return func(ctx *Ctx) Result {
			session, err := ctx.App.Auth.VerifySession(ctx)
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
			return action(ctx)
		}
	}
}
