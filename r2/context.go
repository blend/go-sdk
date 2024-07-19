/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"context"
)

type parameterizedPathKey struct{}

// WithParameterizedPath adds a path with named parameters to a context. Useful for
// outbound request aggregation for metrics and tracing when route parameters
// are involved.
func WithParameterizedPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, parameterizedPathKey{}, path)
}

// GetParameterizedPath gets a path with named parameters off a context. Useful for
// outbound request aggregation for metrics and tracing when route parameters
// are involved. Relies on OptPathParameterized being added to a Request.
func GetParameterizedPath(ctx context.Context) string {
	path, _ := ctx.Value(parameterizedPathKey{}).(string)
	return path
}

type serviceHostNameKey struct{}

// WithServiceHostName adds a service hostName to a context. Useful for
// outbound request aggregation for metrics and tracing when route parameters
// are involved.
func WithServiceHostName(ctx context.Context, serviceHostName string) context.Context {
	return context.WithValue(ctx, serviceHostNameKey{}, serviceHostName)
}

// GetServiceHostName gets a service hostName off a context. Useful for
// outbound request aggregation for metrics and tracing when route parameters
// are involved. Relies on OptServiceHostName being added to a Request.
func GetServiceHostName(ctx context.Context) string {
	serviceHostName, _ := ctx.Value(serviceHostNameKey{}).(string)
	return serviceHostName
}
