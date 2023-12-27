/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package timeutil

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestMinMax(t *testing.T) {
	assert := assert.New(t)
	values := []time.Time{
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -2),
		time.Now().AddDate(0, 0, -3),
		time.Now().AddDate(0, 0, -4),
	}
	min, max := MinMax(values...)
	assert.Equal(values[3], min)
	assert.Equal(values[0], max)
}

func TestMinMaxReversed(t *testing.T) {
	assert := assert.New(t)
	values := []time.Time{
		time.Now().AddDate(0, 0, -4),
		time.Now().AddDate(0, 0, -2),
		time.Now().AddDate(0, 0, -3),
		time.Now().AddDate(0, 0, -1),
	}
	min, max := MinMax(values...)
	assert.Equal(values[0], min)
	assert.Equal(values[3], max)
}

func TestMinMaxEmpty(t *testing.T) {
	assert := assert.New(t)
	values := []time.Time{}
	min, max := MinMax(values...)
	assert.Equal(time.Time{}, min)
	assert.Equal(time.Time{}, max)
}
