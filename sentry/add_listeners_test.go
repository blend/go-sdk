/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sentry

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/configmeta"
	"github.com/blend/go-sdk/logger"
)

func TestAddListeners_Default(t *testing.T) {
	its := assert.New(t)

	its.Nil(AddListeners(nil, configmeta.Meta{}, Config{}))

	log := logger.None()
	its.Nil(AddListeners(log, configmeta.Meta{}, Config{}))
	its.False(log.HasListeners(logger.Error))
	its.False(log.HasListeners(logger.Fatal))

	its.Nil(AddListeners(log, configmeta.Meta{}, Config{DSN: "http://foo@example.org/1"}))
	its.True(log.HasListeners(logger.Error))
	its.True(log.HasListeners(logger.Fatal))
	its.False(log.HasListeners(logger.Warning))

	its.True(log.HasListener(logger.Error, ListenerName))
	its.True(log.HasListener(logger.Fatal, ListenerName))
}

func TestAddListeners_FlagsOption(t *testing.T) {
	its := assert.New(t)

	log := logger.None()
	its.Nil(AddListeners(log, configmeta.Meta{}, Config{DSN: "http://foo@example.org/1"}, AddListenersOptionFlags(logger.Warning, logger.Error)))
	its.True(log.HasListeners(logger.Error))
	its.True(log.HasListeners(logger.Warning))
	its.False(log.HasListeners(logger.Fatal))

	its.True(log.HasListener(logger.Error, ListenerName))
	its.True(log.HasListener(logger.Warning, ListenerName))
}
