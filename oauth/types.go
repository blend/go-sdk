/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
func (e Error) Error() string	{ return string(e) }
