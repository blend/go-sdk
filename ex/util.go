package ex

// ErrClass returns the exception class or the error message.
// This depends on if the err is itself an exception or not.
func ErrClass(err interface{}) string {
	if err == nil {
		return ""
	}
	if ex := As(err); ex != nil && ex.Class != nil {
		return ex.Class.Error()
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
	if ex := As(err); ex != nil && ex.Class != nil {
		return ex.Message
	}
	return ""
}

// Is is a helper function that returns if an error is an ex.
func Is(err interface{}, cause error) bool {
	if err == nil || cause == nil {
		return false
	}
	if typed, isTyped := err.(*Ex); isTyped && typed.Class != nil {
		return (typed.Class == cause) || (typed.Class.Error() == cause.Error())
	}
	if typed, ok := err.(error); ok && typed != nil {
		return (err == cause) || (typed.Error() == cause.Error())
	}
	return err == cause

}

// Inner returns an inner error if the error is an ex.
func Inner(err interface{}) error {
	if typed := As(err); typed != nil {
		return typed.Inner
	}
	return nil
}

// As is a helper method that returns an error as an ex.
func As(err interface{}) *Ex {
	if typed, typedOk := err.(*Ex); typedOk {
		return typed
	}
	return nil
}
