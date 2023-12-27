/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stats

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func TestAddDefaultTagsFromEnv(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set("SERVICE_NAME", "someservice")
	env.Env().Set("SERVICE_ENV", "sandbox")
	env.Env().Set("HOSTNAME", "somecontainer")
	env.Env().Set("CLUSTER_NAME", "somecluster")

	// Handles nil collector
	AddDefaultTagsFromEnv(nil)

	collector := NewMockCollector(32)
	AddDefaultTagsFromEnv(collector)

	tags := collector.DefaultTags()
	assert.Len(tags, 3)
	assert.Equal("service:someservice", tags[0])
	assert.Equal("env:sandbox", tags[1])
	assert.Equal("cluster_name:somecluster", tags[2])
}

func TestAddDefaultTags(t *testing.T) {
	assert := assert.New(t)

	// Handles nil collector
	AddDefaultTagsFromEnv(nil)

	collector := NewMockCollector(32)
	AddDefaultTags(collector, "someservice", "sandbox", "somecluster")

	tags := collector.DefaultTags()
	assert.Len(tags, 3)
	assert.Equal("service:someservice", tags[0])
	assert.Equal("env:sandbox", tags[1])
	assert.Equal("cluster_name:somecluster", tags[2])
}
