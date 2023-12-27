/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package vault

import (
	"context"
)

type clientKey struct{}

// WithClient sets the vault client on a given context.
func WithClient(ctx context.Context, client Client) context.Context {
	return context.WithValue(ctx, clientKey{}, client)
}

// GetClient gets a vault client on a context.
func GetClient(ctx context.Context) Client {
	if value := ctx.Value(clientKey{}); value != nil {
		if typed, ok := value.(Client); ok {
			return typed
		}
	}
	return nil
}
