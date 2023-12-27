/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
