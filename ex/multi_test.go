/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"fmt"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMulti(t *testing.T) {
	it := assert.New(t)

	ex0 := New(New("hi0"))
	ex1 := New(fmt.Errorf("hi1"))
	ex2 := New("hi2")

	m := Append(ex0, ex1, ex2)

	it.True(strings.HasPrefix(m.Error(), `3 errors occurred:`), m.Error()) //todo, make this test more strict

	it.Len(m.(Multi).WrappedErrors(), 3)

	it.NotNil(m.(Multi).Unwrap())
}
