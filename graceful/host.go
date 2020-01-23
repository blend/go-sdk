package graceful

import (
	"os/signal"
	"sync"

	"github.com/blend/go-sdk/ex"
)

// Host gracefully stops a set hosted processes based on a set of variadic options.
// A "Graceful" processes *must* block on start.
// Fatal errors will be returned, that is, errors that are returned by either .Start() or .Stop().
// Panics are not caught by graceful, and it is assumed that your .Start() or .Stop methods will catch relevant panics.
func Host(hosted []Graceful, opts ...Option) error {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	shutdown := make(chan struct{})
	abortWaitShutdown := make(chan struct{})
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
				safely(func() { close(serverExited) }) // close the server exited channel, but do so safely
				waitServerExited.Done()                // signal the normal exit process is done
			}()
			if err := instance.Start(); err != nil {
				errors <- err
			}
			return
		}(hostedInstance)

		// wait to stop the instance
		go func(instance Graceful) {
			defer waitShutdownComplete.Done()
			select {
			case <-shutdown: // tell the hosted process to stop "gracefully"
				if err := instance.Stop(); err != nil {
					errors <- err
				}
				return
			case <-abortWaitShutdown: // a server has exited on its own
				return // clean up this goroutine
			}
		}(hostedInstance)
	}

	select {
	case <-options.ShutdownSignal: // if we've issued a shutdown, wait for the server to exit
		signal.Stop(options.ShutdownSignal)
		close(shutdown)
		waitShutdownComplete.Wait()
		waitServerExited.Wait()

	case <-options.UpdateSignal:
		for _, h := range hosted {
			if typed, ok := h.(Updater); ok {
				go func(u Updater) {
					if err := u.Update(); err != nil {
						if options.UpdateErrors != nil {
							options.UpdateErrors <- err
						}
					}
				}(typed)
			}
		}

	case <-serverExited: // if any of the servers exited on their own
		close(abortWaitShutdown) // quit the signal listener
		waitShutdownComplete.Wait()
	}
	if len(errors) > 0 {
		return <-errors
	}
	return nil
}

func safely(action func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ex.New(r)
		}
	}()
	action()
	return
}
