/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sentry

import "context"

type contextFingerprintKey struct{}

// GetFingerprint gets a context specific fingerprint from the context.
// You can set this with `WithFingerprint(...)`. It will override the default behavior
// of setting the fingerprint to the logger path + err.Error().
func GetFingerprint(ctx context.Context) []string {
	if ctx == nil {
		return nil
	}
	if value := ctx.Value(contextFingerprintKey{}); value != nil {
		if typed, ok := value.([]string); ok {
			return typed
		}
	}
	return nil
}

// WithFingerprint sets the context fingerprint. You can use this to override the default
// fingerprint value submitted by the SDK to sentry.
func WithFingerprint(ctx context.Context, fingerprint ...string) context.Context {
	return context.WithValue(ctx, contextFingerprintKey{}, fingerprint)
}
