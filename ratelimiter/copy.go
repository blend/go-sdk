/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ratelimiter

import (
	"context"
	"errors"
	"io"
	"time"
)

const (
	// DefaultCopyChunkSizeBytes is the write chunk size in bytes.
	DefaultCopyChunkSizeBytes = 32 * 1024
)

// CopyOptions are options for the throttled copy.
type CopyOptions struct {
	RateBytes   int64
	RateQuantum time.Duration
	ChunkSize   int
	Buffer      []byte
	OnWrite     func(int, time.Duration)
}

// CopyOption mutates CopyOptions.
type CopyOption func(*CopyOptions)

// OptCopyRateBytes sets the bytes portion of the rate.
func OptCopyRateBytes(b int64) CopyOption {
	return func(o *CopyOptions) { o.RateBytes = b }
}

// OptCopyRateQuantum sets the quantum portion of the rate.
func OptCopyRateQuantum(q time.Duration) CopyOption {
	return func(o *CopyOptions) { o.RateQuantum = q }
}

// OptCopyChunkSize sets the quantum portion of the rate.
func OptCopyChunkSize(cs int) CopyOption {
	return func(o *CopyOptions) { o.ChunkSize = cs }
}

// OptCopyBuffer sets the buffer for the copy.
func OptCopyBuffer(buf []byte) CopyOption {
	return func(o *CopyOptions) { o.Buffer = buf }
}

// OptCopyOnWrite sets the on write handler for the copy.
func OptCopyOnWrite(handler func(bytesWritten int, elapsed time.Duration)) CopyOption {
	return func(o *CopyOptions) { o.OnWrite = handler }
}

// errCopyInvalidWrite means that a write returned an impossible count.
var errCopyInvalidWrite = errors.New("throttled copy; invalid write result")

// errCopyInvalidChunkSize means that the user provided a < 1 chunk size.
var errCopyInvalidChunkSize = errors.New("throttled copy; invalid chunk size")

// errCopyInvalidOnWrite means that the user provided a nil write handler.
var errCopyInvalidOnWrite = errors.New("throttled copy; invalid on write handler")

// Copy copies from the src reader to the dst writer.
func Copy(ctx context.Context, dst io.Writer, src io.Reader, opts ...CopyOption) (written int64, err error) {
	options := CopyOptions{
		RateBytes:   10 * (1 << 27), // 10gbit in bytes, or (10*(2^30))/8
		RateQuantum: time.Second,
		ChunkSize:   DefaultCopyChunkSizeBytes,
		OnWrite:     func(_ int, _ time.Duration) {},
	}
	for _, opt := range opts {
		opt(&options)
	}

	if options.ChunkSize <= 0 {
		err = errCopyInvalidChunkSize
		return
	}
	if options.OnWrite == nil {
		err = errCopyInvalidOnWrite
		return
	}

	if options.Buffer == nil {
		size := options.ChunkSize
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		options.Buffer = make([]byte, size)
	}

	var nr, nw int
	var er, ew error
	var ts time.Time
	wait := Wait{
		NumberOfActions: options.RateBytes,
		Quantum:         options.RateQuantum,
	}
	var after *time.Timer
	for {
		ts = time.Now()
		nr, er = src.Read(options.Buffer)
		if nr > 0 {
			nw, ew = dst.Write(options.Buffer[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errCopyInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		if err = wait.WaitTimer(ctx, int64(nw), time.Since(ts), after); err != nil {
			return
		}
		options.OnWrite(nw, time.Since(ts))
	}
	return written, err
}
