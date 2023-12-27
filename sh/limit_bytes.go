/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sh

import (
	"io"
	"sync"

	"github.com/blend/go-sdk/ex"
)

// Errors
const (
	ErrMaxBytesWriterCapacityLimit ex.Class = "write failed; maximum capacity reached or would be exceede"
)

// LimitBytes returns a new max bytes writer.
func LimitBytes(maxBytes int, inner io.Writer) *MaxBytesWriter {
	return &MaxBytesWriter{
		max:	maxBytes,
		inner:	inner,
	}
}

// MaxBytesWriter returns a maximum bytes writer.
type MaxBytesWriter struct {
	sync.Mutex

	max	int
	count	int
	inner	io.Writer
}

// Max returns the maximum bytes allowed.
func (mbw *MaxBytesWriter) Max() int {
	return int(mbw.max)
}

// Count returns the current written bytes..
func (mbw *MaxBytesWriter) Count() int {
	return int(mbw.count)
}

func (mbw *MaxBytesWriter) Write(contents []byte) (int, error) {
	mbw.Lock()
	defer mbw.Unlock()

	if mbw.count >= mbw.max {
		return 0, ex.New(ErrMaxBytesWriterCapacityLimit)
	}
	if len(contents)+mbw.count >= mbw.max {
		return 0, ex.New(ErrMaxBytesWriterCapacityLimit)
	}

	written, err := mbw.inner.Write(contents)
	mbw.count += written
	if err != nil {
		return written, ex.New(err)
	}
	return written, nil
}
