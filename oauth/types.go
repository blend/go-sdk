/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package oauth

// Any is a loose type alias to interface{}
type Any = interface{}

// Labels is a loose type alias to map[string]string
type Labels = map[string]string

// Values is a loose type alias to map[string]interface{}
type Values = map[string]interface{}

// Error is an error string.
type Error string

// Error returns the error as a string.
func (e Error) Error() string { return string(e) }
