/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package assert

import "context"

type testNameKey struct{}

// WithTestName sets the test name.
func WithTestName(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, testNameKey{}, id)
}

// GetTestName gets the test name for a test run context.
func GetTestName(ctx context.Context) string {
	if value := ctx.Value(testNameKey{}); value != nil {
		if typed, ok := value.(string); ok {
			return typed
		}
	}
	return ""
}

type contextIDKey struct{}

// WithContextID sets the test context id.
func WithContextID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextIDKey{}, id)
}

// GetContextID gets the context id for a test run.
func GetContextID(ctx context.Context) string {
	if value := ctx.Value(contextIDKey{}); value != nil {
		if typed, ok := value.(string); ok {
			return typed
		}
	}
	return ""
}
