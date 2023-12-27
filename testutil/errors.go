/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package testutil

// NOTE: Ensure that
//   - `SingleError` satisfies `ErrorProducer`.
//   - `SliceErrors` satisfies `ErrorProducer`.
var (
	_	ErrorProducer	= (*SingleError)(nil)
	_	ErrorProducer	= (*SliceErrors)(nil)
)

// ErrorProducer is an interface that defines an error factory.
type ErrorProducer interface {
	NextError() error
}

// SingleError satisfies ErrorProducer for a single error.
type SingleError struct {
	Error error
}

// NextError produces the "next" error.
func (se *SingleError) NextError() error {
	return se.Error
}

// SliceErrors satisfies ErrorProducer for a slice of errors.
type SliceErrors struct {
	Errors	[]error
	Index	int
}

// NextError produces the "next" error. (This is not concurrency safe.)
func (se *SliceErrors) NextError() error {
	index := se.Index
	se.Index++

	if index < 0 {
		return nil
	}

	if index >= len(se.Errors) {
		return nil
	}

	return se.Errors[index]
}
