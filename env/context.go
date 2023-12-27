/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package env

import (
	"context"
)

type varsKey struct{}

// WithVars adds environment variables to a context.
func WithVars(ctx context.Context, vars Vars) context.Context {
	return context.WithValue(ctx, varsKey{}, vars)
}

// GetVars gets environment variables from a context.
func GetVars(ctx context.Context) Vars {
	if raw := ctx.Value(varsKey{}); raw != nil {
		if typed, ok := raw.(Vars); ok {
			return typed
		}
	}
	return Env()
}
