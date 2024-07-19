/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sanitize

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_PathUUIDs(t *testing.T) {
	testCases := [...]struct {
		Input    string
		Expected string
	}{
		{Input: "", Expected: ""},
		{Input: "/", Expected: "/"},
		{Input: "/foo", Expected: "/foo"},
		{Input: "/foo/", Expected: "/foo/"},
		{Input: "//foo", Expected: "//foo"},
		{Input: "//foo//", Expected: "//foo//"},
		{Input: "foo", Expected: "foo"},
		{Input: "foo/", Expected: "foo/"},
		{Input: "foo//", Expected: "foo//"},

		{Input: "/foo/ffbb4178-1ef1-11ec-9255-00155d4fd2f2", Expected: "/foo/?"},
		{Input: "foo/ffbb4178-1ef1-11ec-9255-00155d4fd2f2", Expected: "foo/?"},
		{Input: "/foo/ffbb41781ef111ec925500155d4fd2f2", Expected: "/foo/?"},
		{Input: "foo/ffbb41781ef111ec925500155d4fd2f2", Expected: "foo/?"},

		{Input: "/foo/ffbb417", Expected: "/foo/ffbb417"},
		{Input: "foo/ffbb417", Expected: "foo/ffbb417"},
		{Input: "/foo/ffbb417/", Expected: "/foo/ffbb417/"},
		{Input: "foo/ffbb417/", Expected: "foo/ffbb417/"},

		{Input: "/foo/ffbb4178-1ef1-11ec-9255-00155d4fd2f2/", Expected: "/foo/?/"},
		{Input: "foo/ffbb4178-1ef1-11ec-9255-00155d4fd2f2/", Expected: "foo/?/"},
		{Input: "/foo/ffbb41781ef111ec925500155d4fd2f2/", Expected: "/foo/?/"},
		{Input: "foo/ffbb41781ef111ec925500155d4fd2f2/", Expected: "foo/?/"},

		{Input: "/foo/ffbb4178-1ef1-11ec-9255-00155d4fd2f2/bar", Expected: "/foo/?/bar"},
		{Input: "foo/ffbb4178-1ef1-11ec-9255-00155d4fd2f2/bar", Expected: "foo/?/bar"},
		{Input: "/foo/ffbb4178-1ef1-11ec-9255-00155d4fd2f2/bar/0db8cb24-1ef2-11ec-bc9c-00155d4fd2f2", Expected: "/foo/?/bar/?"},
		{Input: "foo/ffbb4178-1ef1-11ec-9255-00155d4fd2f2/bar/0db8cb24-1ef2-11ec-bc9c-00155d4fd2f2", Expected: "foo/?/bar/?"},
		{Input: "/foo/ffbb4178-1ef1-11ec-9255-00155d4fd2f2/bar/0db8cb24-1ef2-11ec-bc9c-00155d4fd2f2/", Expected: "/foo/?/bar/?/"},
		{Input: "foo/ffbb4178-1ef1-11ec-9255-00155d4fd2f2/bar/0db8cb24-1ef2-11ec-bc9c-00155d4fd2f2/", Expected: "foo/?/bar/?/"},
		{Input: "/foo/ffbb4178-1ef1-11ec-9255-00155d4fd2f2/bar/0db8cb24-1ef2-11ec-bc9c-00155d4fd2f2/baz", Expected: "/foo/?/bar/?/baz"},
		{Input: "foo/ffbb4178-1ef1-11ec-9255-00155d4fd2f2/bar/0db8cb24-1ef2-11ec-bc9c-00155d4fd2f2/baz", Expected: "foo/?/bar/?/baz"},
	}

	for _, tc := range testCases {
		t.Run(tc.Input, func(t2 *testing.T) {
			its := assert.New(t2)
			its.Equal(tc.Expected, PathUUIDs(tc.Input))
		})
	}
}
