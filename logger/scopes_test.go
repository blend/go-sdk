/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Scopes_IsEnabled_all(t *testing.T) {
	its := assert.New(t)

	all := ScopesAll()
	its.True(all.IsEnabled())
	its.True(all.IsEnabled("test0"))
	its.True(all.IsEnabled("test0", "test1"))
}

func Test_Scopes_IsEnabled_all_explicit(t *testing.T) {
	its := assert.New(t)

	all := NewScopes(ScopeAll)
	its.True(all.All())
	its.True(all.IsEnabled())
	its.True(all.IsEnabled("test0"))
	its.True(all.IsEnabled("test0", "test1"))
}

func Test_Scopes_IsEnabled_none(t *testing.T) {
	its := assert.New(t)

	none := ScopesNone()
	its.True(none.None())
	its.False(none.IsEnabled())
	its.False(none.IsEnabled("test0"))
	its.False(none.IsEnabled("test0", "test1"))
}

func Test_Scopes_SetNone(t *testing.T) {
	its := assert.New(t)

	none := NewScopes("foo", "bar")
	its.True(none.IsEnabled("foo"))
	its.True(none.IsEnabled("bar"))
	none.SetNone()
	its.True(none.None())
	its.False(none.IsEnabled())
	its.False(none.IsEnabled("foo"))
	its.False(none.IsEnabled("bar"))
}

func Test_Scopes_IsEnabled_strict(t *testing.T) {
	its := assert.New(t)

	scopes := NewScopes(
		"test0/test1",
		"foo0/foo1",
	)
	its.True(scopes.IsEnabled())

	its.False(scopes.all)

	its.False(scopes.IsEnabled("test0"))
	its.True(scopes.IsEnabled("test0", "test1"))
	its.False(scopes.IsEnabled("test0", "test2"))
	its.False(scopes.IsEnabled("test2", "test1"))

	its.False(scopes.IsEnabled("foo0"))
	its.True(scopes.IsEnabled("foo0", "foo1"))
	its.False(scopes.IsEnabled("foo0", "foo2"))
	its.False(scopes.IsEnabled("foo2", "foo1"))
}

func Test_Scopes_IsEnabled_strict_explicitDisable(t *testing.T) {
	its := assert.New(t)

	all := NewScopes(
		"test0/test1",
		"-foo0/foo1",
	)
	its.True(all.IsEnabled())

	its.False(all.IsEnabled("test0"))
	its.True(all.IsEnabled("test0", "test1"))
	its.False(all.IsEnabled("test0", "test2"))
	its.False(all.IsEnabled("test2", "test1"))

	its.False(all.IsEnabled("foo0"))
	its.False(all.IsEnabled("foo0", "foo1"))
	its.False(all.IsEnabled("foo0", "foo2"))
	its.False(all.IsEnabled("foo2", "foo1"))
}

func Test_Scopes_IsEnabled_wildcard(t *testing.T) {
	its := assert.New(t)

	scopes := NewScopes(
		"test0/*",
		"foo0/foo1",
	)
	its.True(scopes.IsEnabled())

	its.False(scopes.IsEnabled("test0"))
	its.True(scopes.IsEnabled("test0", "test1"))
	its.True(scopes.IsEnabled("test0", "test2"))
	its.False(scopes.IsEnabled("test1", "test0"))

	its.False(scopes.IsEnabled("foo0"))
	its.True(scopes.IsEnabled("foo0", "foo1"))
	its.False(scopes.IsEnabled("foo0", "foo2"))
	its.False(scopes.IsEnabled("foo2", "foo1"))
}

func Test_Scopes_IsEnabled_inputAll_strict_explicitDisable(t *testing.T) {
	its := assert.New(t)

	all := NewScopes(
		ScopeAll,
		"-foo0/foo1",
	)
	its.True(all.IsEnabled())

	its.True(all.IsEnabled("test0"))
	its.True(all.IsEnabled("test0", "test1"))
	its.True(all.IsEnabled("test0", "test2"))
	its.True(all.IsEnabled("test2", "test1"))

	its.True(all.IsEnabled("foo0")) // doesn't match glob
	its.False(all.IsEnabled("foo0", "foo1"))
	its.True(all.IsEnabled("foo0", "foo2"))
	its.True(all.IsEnabled("foo2", "foo1"))
}

func Test_Scopes_IsEnabled_inputAll_wildcard_explicitDisable(t *testing.T) {
	its := assert.New(t)

	all := NewScopes(
		ScopeAll,
		"-foo0/*",
	)
	its.True(all.IsEnabled())

	its.True(all.IsEnabled("test0"))
	its.True(all.IsEnabled("test0", "test1"))
	its.True(all.IsEnabled("test0", "test2"))
	its.True(all.IsEnabled("test2", "test1"))

	its.True(all.IsEnabled("foo0")) // doesn't match glob
	its.False(all.IsEnabled("foo0", "foo1"))
	its.False(all.IsEnabled("foo0", "foo2"))
	its.True(all.IsEnabled("foo2", "foo1"))
}

func Test_Scopes_IsEnabled_wildcard_explicitDisable(t *testing.T) {
	its := assert.New(t)

	scopes := NewScopes(
		"-test0/*",
		"foo0/foo1",
	)
	its.True(scopes.IsEnabled())

	its.False(scopes.IsEnabled("test0"))
	its.False(scopes.IsEnabled("test0", "test1"))
	its.False(scopes.IsEnabled("test0", "test2"))
	its.False(scopes.IsEnabled("test1", "test0"))

	its.False(scopes.IsEnabled("foo0"))
	its.True(scopes.IsEnabled("foo0", "foo1"))
	its.False(scopes.IsEnabled("foo0", "foo2"))
	its.False(scopes.IsEnabled("foo2", "foo1"))
}

func Test_Scopes_String(t *testing.T) {
	its := assert.New(t)

	its.Equal("*", ScopesAll().String())
	its.Equal("*", NewScopes(ScopeAll, "foo/bar/*", "test/testo").String())
	its.Equal("*, -test/testo", NewScopes(ScopeAll, "foo/bar/*", "-test/testo").String())
}

func Test_Scopes_isExplicitlyDisabled(t *testing.T) {
	its := assert.New(t)

	s := NewScopes("*", "-foo/*")

	its.False(s.isScopeExplicitlyDisabled("bar"))
	its.False(s.isScopeExplicitlyDisabled("foo"))
	its.False(s.isScopeExplicitlyDisabled("bar/foo"))
	its.True(s.isScopeExplicitlyDisabled("foo/bar"))
}

func Test_Scopes_matches(t *testing.T) {
	its := assert.New(t)

	its.True(NewScopes("*").matches("foo/bar", "foo/*"))
	its.False(NewScopes("*").matches("foo/*", "foo/bar"))
}
