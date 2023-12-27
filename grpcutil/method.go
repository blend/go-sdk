/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcutil

import "context"

type methodKey struct{}

// WithMethod adds a method to a context as a value.
func WithMethod(ctx context.Context, fullMethod string) context.Context {
	return context.WithValue(ctx, methodKey{}, fullMethod)
}

// GetMethod returns the rpc method from the context.
func GetMethod(ctx context.Context) string {
	if typed, ok := ctx.Value(methodKey{}).(string); ok {
		return typed
	}
	return ""
}
