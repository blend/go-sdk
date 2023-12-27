/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

// ErrInner returns an inner error if the error is an ex.
func ErrInner(err interface{}) error {
	if typed := As(err); typed != nil {
		return typed.Inner
	}
	if typed, ok := err.(InnerProvider); ok && typed != nil {
		return typed.Inner()
	}
	return nil
}
