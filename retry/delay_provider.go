/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package retry

import (
	"context"
	"time"
)

// DelayProvider is a provider for retry delays.
type DelayProvider func(context.Context, uint) time.Duration

// ConstantDelay returns a constant delay provider.
func ConstantDelay(d time.Duration) DelayProvider {
	return func(_ context.Context, _ uint) time.Duration {
		return d
	}
}

// ExponentialBackoff is a backoff provider that doubles the base delay each attempt.
func ExponentialBackoff(d time.Duration) DelayProvider {
	return func(_ context.Context, attempt uint) time.Duration {
		return d * (1 << attempt)
	}
}
