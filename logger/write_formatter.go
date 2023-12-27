/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
