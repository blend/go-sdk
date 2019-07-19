package validate

// Validated is a root interface for validated types.
type Validated interface {
	Validate() error
}

// IsValidated returns if an object implements Validated.
// It does not detect if the Validate function has been called.
func IsValidated(obj interface{}) bool {
	_, ok := obj.(Validated)
	return ok
}

// AsValidated returns the object as a validated and a no-op if it's not.
func AsValidated(obj interface{}) Validated {
	typed, ok := obj.(Validated)
	if !ok {
		return NoOp{}
	}
	return typed
}
