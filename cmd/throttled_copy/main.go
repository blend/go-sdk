/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/blend/go-sdk/fileutil"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/ratelimiter"
)

var (
	rateBytes	= flag.String("rate-bytes", "1024kb", "The throttle rate in bytes")
	rateQuantum	= flag.Duration("rate-quantum", time.Second, "The throttle quantum as a duration")
	verbose		= flag.Bool("verbose", false, "If we should show verbose output")
)

func init() {
	flag.Usage = func() {
		fmt.Printf("throttled_copy SRC DST [flags]\n\nflags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	var src, dst string
	if numArgs := len(flag.Args()); numArgs != 2 {
		flag.Usage()
		os.Exit(1)
	}
	src, dst = flag.Args()[0], flag.Args()[1]
	maybeFatal(throttledCopy(graceful.Background(), src, dst))
}

func throttledCopy(ctx context.Context, src, dst string) error {
	rateBytesValue, err := fileutil.ParseFileSize(*rateBytes)
	if err != nil {
		return err
	}

	var sr io.Reader
	var dw io.Writer
	if src == "-" {
		sr = os.Stdin
	} else {
		sf, err := os.Open(src)
		if err != nil {
			return err
		}
		defer sf.Close()
		sr = sf
	}
	if dst == "-" {
		dw = os.Stdout
	} else {
		df, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer df.Close()
		dw = df
	}

	opts := []ratelimiter.CopyOption{
		ratelimiter.OptCopyRateBytes(rateBytesValue),
		ratelimiter.OptCopyRateQuantum(*rateQuantum),
	}

	if *verbose {
		var written int64
		opts = append(opts,
			ratelimiter.OptCopyOnWrite(func(wr int, e time.Duration) {
				written += int64(wr)
				targetRate := float64(rateBytesValue) / (float64(*rateQuantum) / float64(time.Second))
				rate := float64(wr) / (float64(e) / float64(time.Second))
				fmt.Printf("written: %v, target: %0.2f/s, last: %0.2f/s\n", fileutil.FormatFileSize(written), targetRate, rate)
			}),
		)
	}
	_, err = ratelimiter.Copy(ctx, dw, sr, opts...)
	return err
}

func maybeFatal(err error) {
	if err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
