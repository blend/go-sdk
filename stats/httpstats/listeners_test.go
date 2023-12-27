/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package httpstats

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/webutil"
)

func TestAddListeners(t *testing.T) {
	assert := assert.New(t)

	log := logger.None()
	AddListeners(nil, nil)
	assert.False(log.HasListener(webutil.FlagHTTPRequest, stats.ListenerNameStats))
	AddListeners(log, stats.NewMockCollector(32))
	assert.True(log.HasListener(webutil.FlagHTTPRequest, stats.ListenerNameStats))
}
