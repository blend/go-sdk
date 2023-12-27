/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

// LazyString returns a StringSource for a given string pointer.
//
// LazyString differs from StringPtr in that it treats empty strings as unset.
// If an empty string is a valid value, use a StringPtr.
func LazyString(value *string) LazyStringSource {
	return LazyStringSource{Value: value}
}

var (
	_ StringSource = (*LazyStringSource)(nil)
)

// LazyStringSource implements the LazyString resolver.
type LazyStringSource struct {
	Value *string
}

// String yields the underlying pointer if references a non-empty string.
func (s LazyStringSource) String(_ context.Context) (*string, error) {
	if s.Value == nil || *s.Value == "" {
		return nil, nil
	}
	return s.Value, nil
}
