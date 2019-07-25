package ex

import "encoding/json"

// ClassProvider is a type that can return an exception class.
type ClassProvider interface {
	Class() error
}

// Class is a string wrapper that implements `error`.
// Use this to implement constant exception causes.
type Class string

// Class implements `error`.
func (c Class) Error() string {
	return string(c)
}

// MarshalJSON implements json.Marshaler.
func (c Class) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}
