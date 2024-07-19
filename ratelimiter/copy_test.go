/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ratelimiter

import (
	"bytes"
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Copy(t *testing.T) {
	its := assert.New(t)

	src := bytes.NewBufferString("this is a test")
	dst := new(bytes.Buffer)

	n, err := Copy(context.Background(), dst, src)
	its.Nil(err)
	its.Equal(14, n)
	its.Equal("this is a test", dst.String())
}
