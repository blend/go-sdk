/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcstats

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/grpcutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
)

func TestAddListeners(t *testing.T) {
	assert := assert.New(t)

	log := logger.None()
	AddListeners(nil, nil)
	assert.False(log.HasListener(grpcutil.FlagRPC, stats.ListenerNameStats))
	AddListeners(log, stats.NewMockCollector(32))
	assert.True(log.HasListener(grpcutil.FlagRPC, stats.ListenerNameStats))
}
