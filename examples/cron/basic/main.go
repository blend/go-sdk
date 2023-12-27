/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
)

// Variables
var (
	N = 1024
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log := logger.All()
	jm := cron.New(cron.OptLog(log))

	for x := 0; x < N; x++ {
		jm.LoadJobs(
			cron.NewJob(
				cron.OptJobName(fmt.Sprintf("load-test-%d", x)),
				cron.OptJobSchedule(cron.EverySecond()),
				cron.OptJobAction(func(ctx context.Context) error {
					select {
					case <-ctx.Done():
						return context.Canceled
					case <-time.After(500 * time.Millisecond):
						return nil
					}
				}),
			),
		)
	}
	if err := graceful.Shutdown(jm); err != nil {
		log.Fatal(err)
	}
}
