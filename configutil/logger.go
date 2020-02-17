package configutil

// MaybeInfof writes an info message if the logger is set.
func MaybeInfof(log Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Infof(format, args...)
}

// MaybeDebugf writes a debug message if the logger is set.
func MaybeDebugf(log Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Debugf(format, args...)
}

// MaybeErrorf writes an error message if the logger is set.
func MaybeErrorf(log Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Errorf(format, args...)
}

// Logger is a type that can satisfy the configutil logger interface.
type Logger interface {
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
	Errorf(string, ...interface{})
}
