package graceful

// Shutdown racefully stops a set hosted processes based on SIGINT or SIGTERM received from the os.
// It will return any errors returned by Start() that are not caused by shutting down the server.
// A "Graceful" processes *must* block on start.
func Shutdown(hosted ...Graceful) error {
	return ShutdownBySignal(hosted,
		OptDefaultShutdownSignal(),
	)
}
