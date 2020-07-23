package envoyutil

import (
	"github.com/blend/go-sdk/ex"
)

// NOTE: Ensure
//       - `XFCCError` satisfies `error`.
var (
	_ error = (*XFCCError)(nil)
)

// XFCCError contains metadata about an XFCC header that could not
// be parsed. This is intended to be used as the body of a 401 Unauthorized
// response.
type XFCCError struct {
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
func (xe *XFCCError) Error() string {
	return xe.Class.Error()
}
