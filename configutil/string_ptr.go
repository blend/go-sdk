/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

// StringPtr returns a StringSource for a given string pointer.
//
// It differs from LazyString in that you can resolve to an empty string
// with a StringPtr, but a LazyString would treat that as unset.
func StringPtr(value *string) *StringPtrSource {
	return &StringPtrSource{Value: value}
}

var (
	_ StringSource = (*StringPtrSource)(nil)
)

// StringPtrSource implements the StringPtr resolver.
type StringPtrSource struct {
	Value *string
}

// String yields the underlying pointer, which can be an empty string.
func (s *StringPtrSource) String(_ context.Context) (*string, error) {
	return s.Value, nil
}
