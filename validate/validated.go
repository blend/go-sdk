package validate

// Validated is a root interface for validated types.
type Validated interface {
	Validate() error
}

// Is returns if an object is validated.
func Is(obj interface{}) bool {
	_, ok := obj.(Validated)
	return ok
}

// As returns the object as a validated and a no-op if it's not.
func As(obj interface{}) Validated {
	typed, ok := obj.(Validated)
	if !ok {
		return NoOp{}
	}
	return typed
}
