/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package env

import "fmt"

// ErrNotFound is returned when a field is not found.
type ErrNotFound struct {
	Key string
}

// Error returns the error's text.
func (e ErrNotFound) Error() string {
	return fmt.Sprintf("value for `%s` was not found.", e.Key)
}
