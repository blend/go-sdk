/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sanitize

// PathSanitizer is a type that can sanitize a url path.
type PathSanitizer interface {
	SanitizePath(path string) string
}

// PathSanitizerFunc implements PathSanitizer.
type PathSanitizerFunc func(path string) string

// SanitizePath implements PathSanitizer.
func (psf PathSanitizerFunc) SanitizePath(path string) string {
	return psf(path)
}

// DefaultPathSanitizerFunc is a default implementation of a path
// sanitizer func that just returns the original path.
func DefaultPathSanitizerFunc(p string) string	{ return p }
