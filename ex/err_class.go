/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

// ErrClass returns the exception class or the error message.
// This depends on if the err is itself an exception or not.
func ErrClass(err interface{}) error {
	if err == nil {
		return nil
	}
	if ex := As(err); ex != nil && ex.Class != nil {
		return ex.Class
	}
	if typed, ok := err.(ClassProvider); ok && typed != nil {
		return typed.Class()
	}
	if typed, ok := err.(error); ok && typed != nil {
		return typed
	}
	return nil
}
