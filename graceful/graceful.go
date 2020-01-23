package graceful

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/blend/go-sdk/ex"
)

// DefaultSignals are the default os signals to capture.
var DefaultSignals = []os.Signal{
	os.Interrupt, syscall.SIGTERM,
}

// Graceful is a server that can start and shutdown.
type Graceful interface {
	// Start the service. This _must_ block.
	Start() error
	// Stop the service.
	Stop() error
}

// Signal returns a channel that listens for a given set of os signals.
func Signal(signals ...os.Signal) chan os.Signal {
	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, signals...)
	return terminateSignal
}

// Shutdown racefully stops a set hosted processes based on SIGINT or SIGTERM received from the os.
// It will return any errors returned by Start() that are not caused by shutting down the server.
// A "Graceful" processes *must* block on start.
func Shutdown(hosted ...Graceful) error {
	return ShutdownBySignal(hosted,
		OptShutdownGracePeriod(0),
		OptRecoverPanics(true),
		OptSignal(Signal(DefaultSignals...)),
	)
}

// OptSignal sets the shutdown signal.
func OptSignal(signal chan os.Signal) ShutdownOption {
	return func(so *ShutdownOptions) { so.Signal = signal }
}

// OptRecoverPanics sets if we should capture panics.
func OptRecoverPanics(recoverPanics bool) ShutdownOption {
	return func(so *ShutdownOptions) { so.RecoverPanics = recoverPanics }
}

// OptShutdownGracePeriod sets a shutdown grace period.
func OptShutdownGracePeriod(gracePeriod time.Duration) ShutdownOption {
	return func(so *ShutdownOptions) { so.ShutdownGracePeriod = gracePeriod }
}

// ShutdownOption is a mutator for shutdown options.
type ShutdownOption func(*ShutdownOptions)

// ShutdownOptions are the options for graceful shutdown.
type ShutdownOptions struct {
	ShutdownGracePeriod time.Duration
	RecoverPanics       bool
	Signal              chan os.Signal
}

// ShutdownBySignal gracefully stops a set hosted processes based on an os signal channel.
// A "Graceful" processes *must* block on start.
func ShutdownBySignal(hosted []Graceful, opts ...ShutdownOption) error {
	var options ShutdownOptions
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
		go func(instance Graceful) {
			defer func() {
				safely(func() { close(serverExited) }) // close the emergency crash channel, but do so safely
				waitServerExited.Done()                // signal the normal exit process is done
			}()

			func() { // if we have to recover panics, do so in a closure.
				defer func() {
					if options.RecoverPanics {
						errors <- ex.New(recover())
					}
				}()

				// `hosted.Start()` should block here.
				if err := instance.Start(); err != nil {
					errors <- err
				}
			}()
			return
		}(hostedInstance)

		go func(instance Graceful) {
			defer waitShutdownComplete.Done()

			select {
			case <-shutdown: // tell the hosted process to stop "gracefully"
				if options.ShutdownGracePeriod > 0 {
					go func() {
						if err := instance.Stop(); err != nil {
							errors <- err
						}
					}()
					select {
					case <-time.After(options.ShutdownGracePeriod):
						break
					case <-abortWaitShutdown:
						return
					}
				} else {
					if err := instance.Stop(); err != nil {
						errors <- err
					}
				}
				return
			case <-abortWaitShutdown: // a server has exited on its own
				return // clean up this goroutine
			}
		}(hostedInstance)
	}

	select {
	case <-options.Signal: // if we've issued a shutdown, wait for the server to exit
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
