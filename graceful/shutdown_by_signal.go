package graceful

import (
	"os/signal"
	"sync"

	"github.com/blend/go-sdk/ex"
)

// ShutdownBySignal gracefully stops a set hosted processes based on a set of variadic options.
// A "Graceful" processes *must* block on start.
// Fatal errors will be returned, that is, errors that are returned by either .Start() or .Stop().
// Panics are not caught by graceful, and it is assumed that your .Start() or .Stop methods will catch relevant panics.
func ShutdownBySignal(hosted []Graceful, opts ...ShutdownOption) error {
	var options ShutdownOptions
	for _, opt := range opts {
		opt(&options)
	}

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
	case <-options.ShutdownSignal: // if we've issued a shutdown, wait for the server to exit
		signal.Stop(options.ShutdownSignal) // unhook the process signal redirects, the next ^c will crash the process etc.
		close(shouldShutdown)
		waitShutdownComplete.Wait()
		waitServerExited.Wait()

	case <-serverExited: // if any of the servers exited on their own
		close(shouldShutdown) // quit the signal listener
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
