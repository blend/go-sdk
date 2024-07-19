/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNestMiddleware(t *testing.T) {
	assert := assert.New(t)

	values := make(chan string, 4)
	createMiddleware := func(v string) Middleware {
		return func(action Action) Action {
			return func(r *Ctx) Result {
				values <- v
				return action(r)
			}
		}
	}

	set0 := []Middleware{
		createMiddleware("set0-0"),
		createMiddleware("set0-1"),
		createMiddleware("set0-2"),
	}

	action := func(_ *Ctx) Result {
		values <- "action"
		return nil
	}

	finalAction := NestMiddleware(action, set0...)
	assert.NotNil(finalAction)
	result := finalAction(nil)
	assert.Nil(result)

	assert.Equal("set0-2", <-values)
	assert.Equal("set0-1", <-values)
	assert.Equal("set0-0", <-values)
	assert.Equal("action", <-values)
}

func TestNestMiddleware_Append(t *testing.T) {
	assert := assert.New(t)

	values := make(chan string, 6)
	createMiddleware := func(v string) Middleware {
		return func(action Action) Action {
			return func(r *Ctx) Result {
				values <- v
				return action(r)
			}
		}
	}

	set0 := []Middleware{
		createMiddleware("set0-0"),
		createMiddleware("set0-1"),
		createMiddleware("set0-2"),
	}
	set1 := []Middleware{
		createMiddleware("set1-0"),
		createMiddleware("set1-1"),
	}

	action := func(_ *Ctx) Result {
		values <- "action"
		return nil
	}

	finalAction := NestMiddleware(action, append(set0, set1...)...)
	assert.NotNil(finalAction)
	result := finalAction(nil)
	assert.Nil(result)

	assert.Equal("set1-1", <-values)
	assert.Equal("set1-0", <-values)
	assert.Equal("set0-2", <-values)
	assert.Equal("set0-1", <-values)
	assert.Equal("set0-0", <-values)
	assert.Equal("action", <-values)
}
