/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package graceful

import (
	"context"
	"sync"

	"github.com/blend/go-sdk/ex"
)

// ShutdownByContext gracefully stops a set hosted processes based on context cancellation.
func ShutdownByContext(ctx context.Context, hosted ...Graceful) error {
	shouldShutdown := make(chan struct{})
	serverExited := make(chan struct{})

	waitShutdownComplete := sync.WaitGroup{}
	waitShutdownComplete.Add(len(hosted))

	waitServerExited := sync.WaitGroup{}
	waitServerExited.Add(len(hosted))

	errors := make(chan error, 2*len(hosted))

	for _, hostedInstance := range hosted {
		// start the instance
		go func(instance Graceful) {
			defer func() {
				_ = safely(func() { close(serverExited) }) // close the server exited channel, but do so safely
				waitServerExited.Done()                    // signal the normal exit process is done
			}()
			if err := instance.Start(); err != nil {
				errors <- err
			}
		}(hostedInstance)

		// wait to stop the instance
		go func(instance Graceful) {
			defer waitShutdownComplete.Done()
			<-shouldShutdown // tell the hosted process to stop "gracefully"
			if err := instance.Stop(); err != nil {
				errors <- err
			}
		}(hostedInstance)
	}

	select {
	case <-ctx.Done(): // if we've issued a shutdown, wait for the server to exit
		close(shouldShutdown)
		waitShutdownComplete.Wait()
		waitServerExited.Wait()
	case <-serverExited: // if any of the servers exited on their own
		close(shouldShutdown) // quit the signal listener
		waitShutdownComplete.Wait()
	}
	if errorCount := len(errors); errorCount > 0 {
		var err error
		for x := 0; x < errorCount; x++ {
			err = ex.Append(err, <-errors)
		}
		return err
	}
	return nil
}
