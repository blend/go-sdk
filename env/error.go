/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
