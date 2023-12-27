/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package status

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_TrackedActionConfig_Resolve(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	cfg := new(TrackedActionConfig)

	err := cfg.Resolve(context.Background())
	its.Nil(err)

	its.Equal(DefaultRedRequestCount, cfg.RedRequestCount)
	its.Equal(DefaultRedRequestPercentage, cfg.RedRequestPercentage)

	its.Equal(DefaultYellowRequestCount, cfg.YellowRequestCount)
	its.Equal(DefaultYellowRequestPercentage, cfg.YellowRequestPercentage)

	its.Equal(DefaultTrackedActionExpiration, cfg.Expiration)
}
