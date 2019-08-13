package ex

import "fmt"

// Option is an exception option.
type Option func(*Ex)

// OptMessage sets the exception message from a given list of arguments with fmt.Sprint(args...).
func OptMessage(args ...interface{}) Option {
	return func(ex *Ex) {
		ex.Message = fmt.Sprint(args...)
	}
}

// OptMessagef sets the exception message from a given list of arguments with fmt.Sprintf(format, args...).
func OptMessagef(format string, args ...interface{}) Option {
	return func(ex *Ex) {
		ex.Message = fmt.Sprintf(format, args...)
	}
}

// OptStackTrace sets the exception stack.
func OptStackTrace(stack StackTrace) Option {
	return func(ex *Ex) {
		ex.StackTrace = stack
	}
}

// OptInner sets an inner or wrapped ex.
func OptInner(inner error) Option {
	return func(ex *Ex) {
		ex.Inner = NewWithStackDepth(inner, DefaultNewStartDepth)
	}
}

// OptInnerClass sets an inner unwrapped exception.
// Use this if you don't want to include a strack trace for a cause.
func OptInnerClass(inner error) Option {
	return func(ex *Ex) {
		ex.Inner = inner
	}
}

// OptLabelValue sets a label value.
func OptLabelValue(key, value string) Option {
	return func(ex *Ex) {
		if ex.Labels == nil {
			ex.Labels = make(map[string]string)
		}
		ex.Labels[key] = value
	}
}

// OptLabels sets the labels on an exception.
func OptLabels(labels map[string]string) Option {
	return func(ex *Ex) {
		ex.Labels = labels
	}
}

// OptAnnotationValue sets an annotation value.
func OptAnnotationValue(key, value string) Option {
	return func(ex *Ex) {
		if ex.Annotations == nil {
			ex.Annotations = make(map[string]string)
		}
		ex.Annotations[key] = value
	}
}

// OptAnnotations sets the annotations on an exception.
func OptAnnotations(annoatations map[string]string) Option {
	return func(ex *Ex) {
		ex.Annotations = annoatations
	}
}
