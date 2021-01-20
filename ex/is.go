package ex

// Is is a helper function that returns if an error is an ex.
//
// It will handle if the err is an exception, a multi-error or a regular error.
// "Isness" is evaluated by if the class of the exception matches the class of the cause
func Is(err interface{}, cause error) bool {
	if err == nil || cause == nil {
		return false
	}
	if typed := As(err); typed != nil {
		if typed.Class == nil {
			return false
		}
		if causeTyped := As(cause); causeTyped != nil {
			if causeTyped.Class == nil {
				return false
			}
			return typed.Class == causeTyped.Class
		}
		return (typed.Class == cause) || (typed.Class.Error() == cause.Error())
	}
	if typed, ok := err.(ClassProvider); ok {
		return typed.Class() == cause || (typed.Class().Error() == cause.Error())
	}

	// handle the case of multi-exceptions
	if multiTyped, ok := err.(Multi); ok {
		for _, multiErr := range multiTyped {
			if Is(multiErr, cause) {
				return true
			}
		}
		return false
	}

	// handle regular errors
	if typed, ok := err.(error); ok && typed != nil {
		return (err == cause) || (typed.Error() == cause.Error())
	}
	// handle ???
	return err == cause
}
