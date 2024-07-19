/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

// Interceptors chains calls to interceptors as a single interceptor.
func Interceptors(interceptors ...Interceptor) Interceptor {
	if len(interceptors) == 0 {
		return nil
	}
	if len(interceptors) == 1 {
		return interceptors[0]
	}

	var curry = func(a, b Interceptor) Interceptor {
		if a == nil && b == nil {
			return nil
		}
		if a == nil {
			return b
		}
		if b == nil {
			return a
		}
		return InterceptorFunc(func(i Actioner) Actioner {
			return b.Intercept(a.Intercept(i))
		})
	}
	interceptor := interceptors[0]
	for _, next := range interceptors[1:] {
		interceptor = curry(interceptor, next)
	}
	return interceptor
}

// Interceptor returns an actioner for a given actioner.
type Interceptor interface {
	Intercept(action Actioner) Actioner
}

var (
	_ Interceptor = (*InterceptorFunc)(nil)
)

// InterceptorFunc is a function that implements action.
type InterceptorFunc func(Actioner) Actioner

// Intercept implements Interceptor for the function.
func (fn InterceptorFunc) Intercept(action Actioner) Actioner {
	return fn(action)
}
