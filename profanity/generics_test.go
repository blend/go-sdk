/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNoGenericDecls(t *testing.T) {
	for _, tc := range []struct {
		Name		string
		Enabled		bool
		Filename	string
		Contents	string
		ErrLine		int
	}{
		{
			Name:		"no generics",
			Enabled:	true,
			Filename:	"main.go",
			Contents: `package main

func SumInts(m map[string]int) int {
	var s int
	for _, v := range m {
	    s += v
	}
	return s
}
`,
			ErrLine:	0,
		},
		{
			Name:		"generic function",
			Enabled:	true,
			Filename:	"main.go",
			Contents: `package main

func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
	    s += v
	}
	return s
}
`,
			ErrLine:	3,
		},
		{
			Name:		"generic function when disabled",
			Enabled:	false,
			Filename:	"main.go",
			Contents: `package main

func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
	    s += v
	}
	return s
}
`,
			ErrLine:	0,
		},
		{
			Name:		"generic type",
			Enabled:	true,
			Filename:	"main.go",
			Contents: `package main

type (
	Foo[T any] struct {
		Bar T
	}
)
`,
			ErrLine:	4,
		},
		{
			Name:		"generic type when disabled",
			Enabled:	false,
			Filename:	"main.go",
			Contents: `package main

type (
	Foo[T any] struct {
		Bar T
	}
)
`,
			ErrLine:	0,
		},
		{
			Name:		"interface union",
			Enabled:	true,
			Filename:	"main.go",
			Contents: `package main

type (
	Num interface {
		int | float64
	}
)
`,
			ErrLine:	5,
		},
		{
			Name:		"underlying type in interface",
			Enabled:	true,
			Filename:	"main.go",
			Contents: `package main

type (
	Int interface {
		~int
	}
)
`,
			ErrLine:	5,
		},
		{
			Name:		"bitwise ops ok outside of interface",
			Enabled:	true,
			Filename:	"main.go",
			Contents: `package main

func Foo() {
	_ = 0 | 1
}
`,
			ErrLine:	0,
		},
		{
			Name:		"not go source file",
			Enabled:	true,
			Filename:	"README.txt",
			Contents:	`README`,
			ErrLine:	0,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			assert := assert.New(t)

			rule := NoGenericDecls{
				Enabled: tc.Enabled,
			}

			res := rule.Check(tc.Filename, []byte(tc.Contents))
			if tc.ErrLine > 0 {
				assert.Nil(res.Err)
				assert.False(res.OK)
				assert.Equal(tc.Filename, res.File)
				assert.Equal(tc.ErrLine, res.Line)
			} else {
				assert.Nil(res.Err)
				assert.True(res.OK)
			}
		})
	}
}
