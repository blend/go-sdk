/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

// Middleware is a func that implements middleware
type Middleware func(Action) Action

// NestMiddleware reads the middleware variadic args and organizes the calls
// recursively in the order they appear. I.e. NestMiddleware(inner, third,
// second, first) will call "first", "second", "third", then "inner".
func NestMiddleware(action Action, middleware ...Middleware) Action {
	if len(middleware) == 0 {
		return action
	}

	a := action
	for _, i := range middleware {
		if i != nil {
			a = i(a)
		}
	}
	return a
}
