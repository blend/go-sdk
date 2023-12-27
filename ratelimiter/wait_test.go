/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ratelimiter

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_Wait_Calculate(t *testing.T) {
	its := assert.New(t)

	wait := Wait{2 * 1024, time.Second}.Calculate(
		4*1024,
		time.Second,
	)
	its.Equal(time.Second, wait)

	wait = Wait{4 * 1024, time.Second}.Calculate(2*1024, time.Second)
	its.Equal(-500*time.Millisecond, wait)

	// originally 4 kilobytes a second
	// or 240 kilobytes in a minute (i.e 60 seconds, or 60 * 4)
	wait = Wait{2 * 1024, time.Second}.Calculate(240*1024, time.Minute)
	its.Equal(time.Minute, wait, "THINK ABOUT IT. HOW MANY MINUTES PRODUCED 240kb")

	wait = Wait{2 * 1024, time.Second}.Calculate(3*1024, time.Second)
	its.Equal(500*time.Millisecond, wait)
}
