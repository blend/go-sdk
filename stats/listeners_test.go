/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestAddErrorListeners(t *testing.T) {
	assert := assert.New(t)

	log := logger.None()
	AddErrorListeners(nil, nil)
	assert.False(log.HasListener(logger.Warning, ListenerNameStats))
	assert.False(log.HasListener(logger.Error, ListenerNameStats))
	assert.False(log.HasListener(logger.Fatal, ListenerNameStats))
	AddErrorListeners(log, NewMockCollector(32))
	assert.True(log.HasListener(logger.Warning, ListenerNameStats))
	assert.True(log.HasListener(logger.Error, ListenerNameStats))
	assert.True(log.HasListener(logger.Fatal, ListenerNameStats))
}
