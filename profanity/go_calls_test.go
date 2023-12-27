/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCallsPassing(t *testing.T) {
	it := assert.New(t)

	file := `package main

import "foo/bar"

func doFoo() {
	return
}

func main() {
	thing := make(map[string]string)
	fmt.Println(foo.Bar)
	println(bar.Foo)
	doFoo()
}
`
	rule := GoCalls([]GoCall{
		{
			Package:	"fmt",
			Func:		"Printf",
		},
	})

	res := rule.Check("main.go", []byte(file))
	it.Nil(res.Err)
	it.True(res.OK)
}

func TestCallsPrintln(t *testing.T) {
	it := assert.New(t)

	file := `package main

import "foo/bar"

func doFoo() {
	return
}

func main() {
	thing := make(map[string]string)
	fmt.Println(foo.Bar)
	println(bar.Foo)
	doFoo()
}
`
	rule := GoCalls([]GoCall{
		{
			Package:	"fmt",
			Func:		"Println",
		},
	})

	res := rule.Check("main.go", []byte(file))
	it.Nil(res.Err)
	it.False(res.OK)
	it.Equal("main.go", res.File)
	it.Equal(11, res.Line)
}

func TestCallsEmptyPackage(t *testing.T) {
	it := assert.New(t)

	file := `package main

import "foo/bar"

func doFoo() {
	return
}

func main() {
	thing := make(map[string]string)
	fmt.Println(foo.Bar)
	println(bar.Foo)
	doFoo()
}
`

	rule := GoCalls([]GoCall{
		{
			Func: "println",
		},
	})

	res := rule.Check("main.go", []byte(file))
	it.Nil(res.Err)
	it.False(res.OK)
	it.Equal("main.go", res.File)
	it.Equal(12, res.Line)
}
