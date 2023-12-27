/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package status

import (
	"context"
	"time"

	"github.com/blend/go-sdk/configutil"
)

// TrackedActionConfig is the configuration for the tracker.
type TrackedActionConfig struct {
	YellowRequestCount	int
	YellowRequestPercentage	float64
	RedRequestCount		int
	RedRequestPercentage	float64
	Expiration		time.Duration
}

// Resolve resolves the config.
func (tc *TrackedActionConfig) Resolve(ctx context.Context) error {
	return configutil.Resolve(ctx,
		configutil.SetInt(&tc.YellowRequestCount,
			configutil.Int(tc.YellowRequestCount),
			configutil.Int(DefaultYellowRequestCount),
		),
		configutil.SetFloat64(&tc.YellowRequestPercentage,
			configutil.Float64(tc.YellowRequestPercentage),
			configutil.Float64(DefaultYellowRequestPercentage),
		),
		configutil.SetInt(&tc.RedRequestCount,
			configutil.Int(tc.RedRequestCount),
			configutil.Int(DefaultRedRequestCount),
		),
		configutil.SetFloat64(&tc.RedRequestPercentage,
			configutil.Float64(tc.RedRequestPercentage),
			configutil.Float64(DefaultRedRequestPercentage),
		),
		configutil.SetDuration(&tc.Expiration,
			configutil.Duration(tc.Expiration),
			configutil.Duration(DefaultTrackedActionExpiration),
		),
	)
}

// ExpirationOrDefault returns an expiration or a default.
func (tc TrackedActionConfig) ExpirationOrDefault() time.Duration {
	if tc.Expiration > 0 {
		return tc.Expiration
	}
	return DefaultTrackedActionExpiration
}
