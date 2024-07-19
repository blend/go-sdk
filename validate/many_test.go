/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

type manyTest struct {
	Foo *string
	Bar *string
	Baz *string
}

func TestManyNil(t *testing.T) {
	assert := assert.New(t)

	refStr := func(val string) *string { return &val }

	bad := manyTest{
		Foo: refStr("foo"),
		Bar: refStr("bar"),
		Baz: refStr("baz"),
	}
	assert.NotNil(Many(bad.Foo, bad.Bar, bad.Baz).Nil()())

	maybe := manyTest{
		Bar: refStr("bar"),
	}
	assert.NotNil(Many(maybe.Foo, maybe.Bar, maybe.Baz).Nil()())

	good := manyTest{}
	assert.Nil(Many(good.Foo, good.Bar, good.Baz).Nil()())
}

func TestManyNotNil(t *testing.T) {
	assert := assert.New(t)

	refStr := func(val string) *string { return &val }

	bad := manyTest{}
	assert.NotNil(Many(bad.Foo, bad.Bar, bad.Baz).NotNil()())

	maybe := manyTest{
		Bar: refStr("bar"),
	}
	assert.Nil(Many(maybe.Foo, maybe.Bar, maybe.Baz).NotNil()())

	good := manyTest{
		Foo: refStr("foo"),
		Bar: refStr("bar"),
		Baz: refStr("baz"),
	}
	assert.Nil(Many(good.Foo, good.Bar, good.Baz).NotNil()())
}

func TestManyOneNotNil(t *testing.T) {
	assert := assert.New(t)

	refStr := func(val string) *string { return &val }

	bad := manyTest{}
	assert.NotNil(Many(bad.Foo, bad.Bar, bad.Baz).OneNotNil()())

	maybe := manyTest{
		Bar: refStr("bar"),
	}
	assert.Nil(Many(maybe.Foo, maybe.Bar, maybe.Baz).OneNotNil()())

	good := manyTest{
		Foo: refStr("foo"),
		Bar: refStr("bar"),
		Baz: refStr("baz"),
	}
	assert.NotNil(Many(good.Foo, good.Bar, good.Baz).OneNotNil()())
}
