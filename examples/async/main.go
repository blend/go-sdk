/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/blend/go-sdk/async"
)

// WorkSize is the amount of work to do.
const WorkSize = 1 << 18

func main() {
	work := make(chan interface{}, WorkSize)

	for x := 0; x < WorkSize; x++ {
		work <- fmt.Sprintf("work-%d", x)
	}

	batch := async.NewBatch(work, func(ctx context.Context, work interface{}) error {
		fmt.Printf("%v\n", work)
		return nil
	}, async.OptBatchParallelism(runtime.NumCPU()))

	batch.Process(context.TODO())
}
