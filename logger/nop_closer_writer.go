/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import "io"

// NopCloserWriter doesn't allow the underlying writer to be closed.
type NopCloserWriter struct {
	io.Writer
}

// Close does not close.
func (ncw NopCloserWriter) Close() error	{ return nil }
