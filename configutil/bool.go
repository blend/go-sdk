/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

// BoolSource is a type that can return a value.
type BoolSource interface {
	// Bool should return a bool if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	Bool(context.Context) (*bool, error)
}

var (
	_ BoolSource = (*BoolValue)(nil)
)

// Bool returns a BoolValue for a given value.
func Bool(value *bool) *BoolValue {
	if value == nil {
		return nil
	}
	typed := BoolValue(*value)
	return &typed
}

// BoolValue implements value provider.
type BoolValue bool

// Bool returns the value for a constant.
func (b *BoolValue) Bool(_ context.Context) (*bool, error) {
	if b == nil {
		return nil, nil
	}
	value := *b
	typed := bool(value)
	return &typed, nil
}
