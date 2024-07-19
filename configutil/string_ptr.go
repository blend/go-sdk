/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
