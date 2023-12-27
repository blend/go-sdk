/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
