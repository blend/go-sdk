/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package selector

import (
	"encoding/json"
	"fmt"
)

const (
	errExpectedNonEmptyKey = "expected non-empty key"
)

// Error is a hard alias to string.
type Error string

// Error implements `error`
func (e Error) Error() string {
	return string(e)
}

// MarshalJSON implements json.Marshaler.
func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

// ParseError is a specific parse error.
type ParseError struct {
	Err		error
	Input		string
	Position	int
	Message		string
}

// Class implements ex.ClassProvider.
func (pe ParseError) Class() error {
	return pe.Err
}

// Unwrap implements unwrap.
func (pe ParseError) Unwrap() error {
	return pe.Err
}

// String implements error.
func (pe ParseError) Error() string {
	if pe.Message != "" {
		return fmt.Sprintf("%q:0:%d: %v; %s", pe.Input, pe.Position, pe.Err, pe.Message)
	}
	return fmt.Sprintf("%q:0:%d: %v", pe.Input, pe.Position, pe.Err)
}
