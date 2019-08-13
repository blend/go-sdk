package ex

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

// As is a helper method that returns an error as an ex.
func As(err interface{}) *Ex {
	if typed, typedOk := err.(*Ex); typedOk {
		return typed
	}
	return nil
}

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

// InnerProvider is a type that returns an inner error.
type InnerProvider interface {
	Inner() error
}

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

// ErrLabels returns the labels for an error.
func ErrLabels(err interface{}) map[string]string {
	if typed := As(err); typed != nil {
		return typed.Labels
	}
	return nil
}

// ErrAnnotations returns the labels for an error.
func ErrAnnotations(err interface{}) map[string]string {
	if typed := As(err); typed != nil {
		return typed.Annotations
	}
	return nil
}
