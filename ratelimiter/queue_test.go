/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ratelimiter

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestQueue_Check(t *testing.T) {
	it := assert.New(t)

	rl := NewQueue(5, time.Second) // 5 actions per second

	now := time.Now()

	rl.Now = Clock(now, 0)
	it.False(rl.Check("a"), "first call to `a` should pass")

	rl.Now = Clock(now, 100*time.Millisecond)
	it.False(rl.Check("b"), "first call to `b` should pass")

	rl.Now = Clock(now, 200*time.Millisecond)
	it.False(rl.Check("b"), "second call to `b` should pass")

	rl.Now = Clock(now, 300*time.Millisecond)
	it.False(rl.Check("b"), "third call to `b` should pass")

	rl.Now = Clock(now, 400*time.Millisecond)
	it.False(rl.Check("b"), "fourth call to `b` should pass")

	rl.Now = Clock(now, 500*time.Millisecond)
	it.False(rl.Check("a"), "second call to `a` in 500ms should pass")

	rl.Now = Clock(now, 600*time.Millisecond)
	it.False(rl.Check("a"), "third call to `a` in 600ms should pass")

	rl.Now = Clock(now, 700*time.Millisecond)
	it.False(rl.Check("a"), "fourth call to `a` in 700ms should pass")

	rl.Now = Clock(now, 800*time.Millisecond)
	it.True(rl.Check("a"), "fifth call to `a` in 800ms should fail")

	rl.Now = Clock(now, 2000*time.Millisecond)
	it.False(rl.Check("a"), "first call to `a` after pause should pass")

	rl.Now = Clock(now, 2100*time.Millisecond)
	it.False(rl.Check("b"), "first call to `b` after pause should pass")

	rl.Now = Clock(now, 2200*time.Millisecond)
	it.False(rl.Check("b"), "second call to `b` after pause should pass")

	rl.Now = Clock(now, 2300*time.Millisecond)
	it.False(rl.Check("b"), "third call to `b` after pause should pass")

	rl.Now = Clock(now, 2300*time.Millisecond)
	it.False(rl.Check("b"), "fourth call to `b` after pause should pass")

	rl.Now = Clock(now, 2400*time.Millisecond)
	it.False(rl.Check("a"), "second call to `a` after pause should pass")

	rl.Now = Clock(now, 2500*time.Millisecond)
	it.False(rl.Check("a"), "third call to `a` after pause should pass")

	rl.Now = Clock(now, 2600*time.Millisecond)
	it.False(rl.Check("a"), "fourth call to `a` after pause should pass")

	rl.Now = Clock(now, 2700*time.Millisecond)
	it.True(rl.Check("a"), "fifth call to `a` after pause should fail")
}
