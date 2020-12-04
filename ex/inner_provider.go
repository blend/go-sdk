package ex

// InnerProvider is a type that returns an inner error.
type InnerProvider interface {
	Inner() error
}
