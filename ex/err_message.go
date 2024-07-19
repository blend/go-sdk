/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

// ErrMessage returns the exception message.
// This depends on if the err is itself an exception or not.
// If it is not an exception, this will return empty string.
func ErrMessage(err interface{}) string {
	if err == nil {
		return ""
	}
	if ex := As(err); ex != nil && ex.Class != nil {
		return ex.Message
	}
	if typed, ok := err.(MessageProvider); ok && typed != nil {
		return typed.Message()
	}
	return ""
}
