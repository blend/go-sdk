/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import "github.com/blend/go-sdk/ex"

// Errors
const (
	ErrInvalidSameSite		ex.Class	= "invalid cookie same site string value"
	ErrParameterMissing		ex.Class	= "parameter missing"
	ErrUnauthorized			ex.Class	= "unauthorized"
	ErrInvalidSplitColonInput	ex.Class	= `split colon input string is not of the form "<first>:<second>"`
)

// ErrIsInvalidSameSite returns if an error is `ErrInvalidSameSite`
func ErrIsInvalidSameSite(err error) bool {
	return ex.Is(err, ErrInvalidSameSite)
}

// ErrIsParameterMissing returns if an error is `ErrParameterMissing`
func ErrIsParameterMissing(err error) bool {
	return ex.Is(err, ErrParameterMissing)
}

// ErrIsUnauthorized returns if an error is `ErrUnauthorized`
func ErrIsUnauthorized(err error) bool {
	return ex.Is(err, ErrUnauthorized)
}

// ErrIsInvalidSplitColonInput returns if an error is `ErrInvalidSplitColonInput`
func ErrIsInvalidSplitColonInput(err error) bool {
	return ex.Is(err, ErrInvalidSplitColonInput)
}
