package exception

// ErrClass returns the exception class or the error message.
// This depends on if the err is itself an exception or not.
func ErrClass(err interface{}) string {
	if err == nil {
		return ""
	}
	if ex := As(err); ex != nil && ex.Class() != nil {
		return ex.Class().Error()
	}
	if typed, ok := err.(error); ok && typed != nil {
		return typed.Error()
	}
	return ""
}

// ErrMessage returns the exception message.
// This depends on if the err is itself an exception or not.
// If it is not an exception, this will return empty string.
func ErrMessage(err interface{}) string {
	if err == nil {
		return ""
	}
	if ex := As(err); ex != nil && ex.Class() != nil {
		return ex.Message()
	}
	return ""
}

// ErrStackTrace returns a stack trace from an `Error` type if it is an
// `Exception`. If it is not an `Exception`, this will return `nil`.
func ErrStackTrace(err interface{}) StackTrace {
	if ex := As(err); ex != nil && ex.Stack() != nil {
		return ex.Stack()
	}
	return nil
}

// Is is a helper function that returns if an error is an exception.
func Is(err interface{}, cause error) bool {
	if err == nil || cause == nil {
		return false
	}
	if typed, isTyped := err.(Exception); isTyped && typed.Class() != nil {
		return (typed.Class() == cause) || (typed.Class().Error() == cause.Error())
	}
	if typed, ok := err.(error); ok && typed != nil {
		return (err == cause) || (typed.Error() == cause.Error())
	}
	return err == cause

}

// Inner returns an inner error if the error is an exception.
func Inner(err interface{}) error {
	if typed := As(err); typed != nil {
		return typed.Inner()
	}
	return nil
}

// As is a helper method that returns an error as an exception.
func As(err interface{}) Exception {
	if typed, typedOk := err.(Exception); typedOk {
		return typed
	}
	return nil
}
