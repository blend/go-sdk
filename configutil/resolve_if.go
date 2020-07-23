package configutil

import "context"

// ResolveIf wraps a resolver in a branch.
func ResolveIf(branch bool, resolver ResolveAction) ResolveAction {
	return func(ctx context.Context) error {
		if branch {
			return resolver(ctx)
		}
		return nil
	}
}

// ResolveIfFunc wraps a resolver in a branch returned from a function.
func ResolveIfFunc(branchFunc func(context.Context) bool, resolver ResolveAction) ResolveAction {
	return func(ctx context.Context) error {
		if branchFunc(ctx) {
			return resolver(ctx)
		}
		return nil
	}
}
