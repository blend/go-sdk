/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package migration

import "context"

type suiteKey struct{}

// WithSuite adds a suite as a value to a context.
func WithSuite(ctx context.Context, suite *Suite) context.Context {
	return context.WithValue(ctx, suiteKey{}, suite)
}

// GetContextSuite gets a suite from a context as a value.
func GetContextSuite(ctx context.Context) *Suite {
	value := ctx.Value(suiteKey{})
	if typed, ok := value.(*Suite); ok {
		return typed
	}
	return nil
}

type labelsKey struct{}

// WithLabel adds a label to the context
func WithLabel(ctx context.Context, label string) context.Context {
	return context.WithValue(ctx, labelsKey{}, append(GetContextLabels(ctx), label))
}

// GetContextLabels gets a group from a context as a value.
func GetContextLabels(ctx context.Context) []string {
	value := ctx.Value(labelsKey{})
	if typed, ok := value.([]string); ok {
		return typed
	}
	return nil
}
