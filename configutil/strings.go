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

// StringsSource is a type that can return a value.
type StringsSource interface {
	// Strings should return a string array if the source has a given value.
	// It should return nil if the value is not present.
	// It should return an error if there was a problem fetching the value.
	Strings(context.Context) ([]string, error)
}

var (
	_ StringsSource = (*Strings)(nil)
)

// Strings implements a value provider.
type Strings []string

// Strings returns the value for a constant.
func (s Strings) Strings(_ context.Context) ([]string, error) {
	return []string(s), nil
}
