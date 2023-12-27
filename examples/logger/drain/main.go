/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/blend/go-sdk/logger"
)

func main() {
	log := logger.MustNew(logger.OptAll())

	log.Listen(logger.Info, "randomly_slow", func(ctx context.Context, e logger.Event) {
		if rand.Float64() < 0.2 {
			fmt.Println("randomly slow start")
			time.Sleep(2000 * time.Millisecond)
			fmt.Println("randomly slow stop")
		}
	})

	infoSignal := time.Tick(100 * time.Millisecond)

	done := time.After(10 * time.Second)

	for {
		select {
		case <-infoSignal:
			log.Infof("this is an info event")
		case <-done:
			fmt.Println("draining")
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
				defer cancel()
				log.DrainContext(ctx)
			}()
			fmt.Println("exiting")
			os.Exit(0)
		}
	}
}
