/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package ex

// ErrStackTrace returns the exception stack trace.
// This depends on if the err is itself an exception or not.
func ErrStackTrace(err interface{}) StackTrace {
	if err == nil {
		return nil
	}
	if ex := As(err); ex != nil && ex.StackTrace != nil {
		return ex.StackTrace
	}
	if typed, ok := err.(StackTraceProvider); ok && typed != nil {
		return typed.StackTrace()
	}
	return nil
}
