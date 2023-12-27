/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stats

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestContextWithCollector(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	collector := NewMockCollector(1)

	assert.Equal(collector, GetCollector(WithCollector(ctx, collector)))
}
