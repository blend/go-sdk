/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"context"
	"io"
)

// WriteFormatter is a formatter for writing events to output writers.
type WriteFormatter interface {
	WriteFormat(context.Context, io.Writer, Event) error
}
