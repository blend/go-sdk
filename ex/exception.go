package ex

// Exception is a meta interface for exceptions.
type Exception interface {
	error
	WithMessage(...interface{}) Exception
	WithMessagef(string, ...interface{}) Exception
	WithInner(error) Exception
}
