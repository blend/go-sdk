package graceful

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Graceful is a server that can start and shutdown.
type Graceful interface {
	// Start the service. This _must_ block.
	Start() error
	// Stop the service.
	Stop() error
}

// Shutdown racefully stops a set hosted processes based on SIGINT or SIGTERM received from the os.
// It will return any errors returned by Start() that are not caused by shutting down the server.
// A "Graceful" processes *must* block on start.
func Shutdown(hosted ...Graceful) error {
	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, os.Interrupt, syscall.SIGTERM)
	return ShutdownBySignal(terminateSignal, hosted...)
}

// ShutdownBySignal gracefully stops a set hosted processes based on an os signal channel.
// A "Graceful" processes *must* block on start.
func ShutdownBySignal(shouldShutdown chan os.Signal, hosted ...Graceful) error {
	shutdown := make(chan struct{})
	abortWaitShutdown := make(chan struct{})
	serverExited := make(chan struct{})

	waitShutdownComplete := sync.WaitGroup{}
	waitShutdownComplete.Add(len(hosted))

	waitServerExited := sync.WaitGroup{}
	waitServerExited.Add(len(hosted))

	errors := make(chan error, 2*len(hosted))

	for _, hostedInstance := range hosted {
		// start the hosted instance
		go func(instance Graceful) {
			defer func() {
				safely(func() { close(serverExited) }) // close the emergency crash channel, but do so safely
				waitServerExited.Done()                // signal the normal exit process is done
			}()

			// `hosted.Start()` should block here.
			if err := instance.Start(); err != nil {
				errors <- err
			}
			return
		}(hostedInstance)

		go func(instance Graceful) {
			defer waitShutdownComplete.Done()

			select {
			case <-shutdown:
				// tell the hosted process to stop "gracefully"
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
	case <-shouldShutdown: // if we've issued a shutdown, wait for the server to exit
		close(shutdown)
		waitShutdownComplete.Wait()
		waitServerExited.Wait()
	case <-serverExited: // if any of the servers exited on their own
		close(abortWaitShutdown) // quit the signal listener
		waitShutdownComplete.Wait()
	}
	if len(errors) > 0 {
		return <-errors
	}
	return nil
}

func safely(action func()) {
	defer func() {
		recover()
	}()
	action()
}
