package ex

// ClassProvider is a type that can return an exception class.
type ClassProvider interface {
	Class() error
}
