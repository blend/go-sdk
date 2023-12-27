/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"context"
)

type collectorKey struct{}

// WithCollector adds the collector to a context.
func WithCollector(ctx context.Context, c Collector) context.Context {
	return context.WithValue(ctx, collectorKey{}, c)
}

// GetCollector gets a collector off a context.
func GetCollector(ctx context.Context) Collector {
	if value := ctx.Value(collectorKey{}); value != nil {
		if typed, ok := value.(Collector); ok {
			return typed
		}
	}
	return nil
}
