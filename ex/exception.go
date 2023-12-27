/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

// Exception is a meta interface for exceptions.
type Exception interface {
	error
	WithMessage(...interface{}) Exception
	WithMessagef(string, ...interface{}) Exception
	WithInner(error) Exception
}
