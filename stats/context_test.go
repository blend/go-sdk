/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
