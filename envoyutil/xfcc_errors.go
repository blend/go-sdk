/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package envoyutil

import (
	"github.com/blend/go-sdk/ex"
)

// NOTE: Ensure
//   - `XFCCExtractionError` satisfies `error`.
//   - `XFCCValidationError` satisfies `error`.
//   - `XFCCFatalError` satisfies `error`.
var (
	_ error = (*XFCCExtractionError)(nil)
	_ error = (*XFCCValidationError)(nil)
	_ error = (*XFCCFatalError)(nil)
)

// XFCCExtractionError contains metadata about an XFCC header that could not
// be parsed or extracted. This is intended to be used as the body of a 401
// Unauthorized response.
type XFCCExtractionError struct {
	// Class can be used to uniquely identify the type of the error.
	Class ex.Class `json:"class" xml:"class"`
	// XFCC contains the XFCC header value that could not be parsed or was
	// invalid in some way.
	XFCC string `json:"xfcc,omitempty" xml:"xfcc,omitempty"`
	// Metadata contains extra information relevant to a specific failure.
	Metadata interface{} `json:"metadata,omitempty" xml:"metadata,omitempty"`
}

// Error satisfies the `error` interface. It is intended to be a unique
// identifier for the error.
func (xee *XFCCExtractionError) Error() string {
	return xee.Class.Error()
}

// IsExtractionError is a helper to check if an error is an `*XFCCExtractionError`.
func IsExtractionError(err error) bool {
	_, ok := err.(*XFCCExtractionError)
	return ok
}

// XFCCValidationError contains metadata about an XFCC header that could not
// be parsed or extracted. This is intended to be used as the body of a 401
// Unauthorized response.
type XFCCValidationError struct {
	// Class can be used to uniquely identify the type of the error.
	Class ex.Class `json:"class" xml:"class"`
	// XFCC contains the XFCC header value that could not be parsed or was
	// invalid in some way.
	XFCC string `json:"xfcc,omitempty" xml:"xfcc,omitempty"`
	// Metadata contains extra information relevant to a specific failure.
	Metadata interface{} `json:"metadata,omitempty" xml:"metadata,omitempty"`
}

// Error satisfies the `error` interface. It is intended to be a unique
// identifier for the error.
func (xve *XFCCValidationError) Error() string {
	return xve.Class.Error()
}

// IsValidationError is a helper to check if an error is an `*XFCCValidationError`.
func IsValidationError(err error) bool {
	_, ok := err.(*XFCCValidationError)
	return ok
}

// XFCCFatalError contains metadata about an unrecoverable failure when parsing
// an XFCC header. A "fatal error" should indicate invalid usage of `envoyutil`
// such as providing a `nil` value for a function interface that must be invoked.
type XFCCFatalError struct {
	// Class can be used to uniquely identify the type of the error.
	Class ex.Class `json:"class" xml:"class"`
	// XFCC contains the XFCC header value that could not be parsed or was
	// invalid in some way.
	XFCC string `json:"xfcc,omitempty" xml:"xfcc,omitempty"`
}

// Error satisfies the `error` interface. It is intended to be a unique
// identifier for the error.
func (xfe *XFCCFatalError) Error() string {
	return xfe.Class.Error()
}

// IsFatalError is a helper to check if an error is an `*XFCCFatalError`.
func IsFatalError(err error) bool {
	_, ok := err.(*XFCCFatalError)
	return ok
}
