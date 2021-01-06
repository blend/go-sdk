package web

// ViewProviderAsDefault sets the context.DefaultResultProvider() equal to context.View().
func ViewProviderAsDefault(action Action) Action {
	return func(ctx *Ctx) Result {
		ctx.DefaultProvider = ctx.Views
		return action(ctx)
	}
}

// JSONProviderAsDefault sets the context.DefaultResultProvider() equal to context.JSON().
func JSONProviderAsDefault(action Action) Action {
	return func(ctx *Ctx) Result {
		ctx.DefaultProvider = JSON
		return action(ctx)
	}
}

// XMLProviderAsDefault sets the context.DefaultResultProvider() equal to context.XML().
func XMLProviderAsDefault(action Action) Action {
	return func(ctx *Ctx) Result {
		ctx.DefaultProvider = XML
		return action(ctx)
	}
}

// TextProviderAsDefault sets the context.DefaultResultProvider() equal to context.Text().
func TextProviderAsDefault(action Action) Action {
	return func(ctx *Ctx) Result {
		ctx.DefaultProvider = Text
		return action(ctx)
	}
}
