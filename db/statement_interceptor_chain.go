/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import "context"

// StatementInterceptorChain nests interceptors such that a list of interceptors is
// returned as a single function.
//
// The nesting is such that if this function is provided arguments as `a, b, c` the returned
// function would be in the form `c(b(a(...)))` i.e. the functions would be run in left to right order.
func StatementInterceptorChain(interceptors ...StatementInterceptor) StatementInterceptor {
	if len(interceptors) == 0 {
		return func(ctx context.Context, label, statement string) (string, error) {
			return statement, nil
		}
	}
	if len(interceptors) == 1 {
		return interceptors[0]
	}
	var nest = func(a, b StatementInterceptor) StatementInterceptor {
		if a == nil {
			return b
		}
		if b == nil {
			return a
		}

		return func(ctx context.Context, label, statement string) (string, error) {
			var err error
			statement, err = b(ctx, label, statement)
			if err != nil {
				return statement, err
			}
			return a(ctx, label, statement)
		}
	}

	var outer StatementInterceptor
	for _, step := range interceptors {
		outer = nest(step, outer)
	}
	return outer
}
