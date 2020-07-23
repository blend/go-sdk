package web

// Middleware is a func that implements middleware
type Middleware func(Action) Action

// NestMiddleware reads the middleware variadic args and organizes the calls recursively
// in the order they appear.
func NestMiddleware(action Action, middleware ...Middleware) Action {
	if len(middleware) == 0 {
		return action
	}

	var outer Middleware
	for _, step := range middleware {
		outer = nest(step, outer)
	}
	return outer(action)
}

func nest(a, b Middleware) Middleware {
	if b == nil {
		return a
	}
	return func(inner Action) Action {
		return a(b(inner))
	}
}
