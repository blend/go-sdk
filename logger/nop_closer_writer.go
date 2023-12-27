/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import "io"

// NopCloserWriter doesn't allow the underlying writer to be closed.
type NopCloserWriter struct {
	io.Writer
}

// Close does not close.
func (ncw NopCloserWriter) Close() error { return nil }
