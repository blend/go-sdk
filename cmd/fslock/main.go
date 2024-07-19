/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/blend/go-sdk/filelock"
)

var (
	flagExclusive = flag.Bool("exclusive", false, "If the lock is exclusive or not")
	flagRemove    = flag.Bool("remove", false, "If the lock file should be removed when complete")
	flagWaitFor   = flag.Duration("wait-for", 10*time.Second, "The duration to hold the lock for")
	flagHoldFor   = flag.Duration("hold-for", 10*time.Second, "The duration to hold the lock for")
)

func init() {
	flag.Parse()
}

func main() {
	filePath := flag.Arg(0)
	if filePath == "" {
		fmt.Fprintf(os.Stderr, "Must provide a filepath; see usage for details")
		os.Exit(1)
	}

	mu := filelock.MutexAt(filePath)

	var lockfn func() (func(), error)
	if *flagExclusive {
		lockfn = mu.Lock
	} else {
		lockfn = mu.RLock
	}

	unlock, err := waitAcquired(lockfn, *flagWaitFor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error acquiring lock for file %q; %+v\n", filePath, err)
		os.Exit(1)
	}
	defer func() {
		unlock()
		if *flagRemove {
			fmt.Println("Removing lock file", filePath)
			_ = os.Remove(filePath)
		}
	}()

	fmt.Println("Acquired lock", filePath, "holding for", *flagHoldFor)
	<-time.After(*flagHoldFor)
}

func waitAcquired(lockfn func() (func(), error), waitFor time.Duration) (unlock func(), err error) {
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		unlock, err = lockfn()
	}()
	select {
	case <-finished:
		return
	case <-time.After(waitFor):
		err = context.Canceled
		return
	}
}
