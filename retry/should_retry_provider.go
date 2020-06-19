package retry

// ShouldRetryProvider is a function that returns if we should retry
// on an error or abort retries.
// Return `true` to continue to retry, and `false` otherwise to abort retries.
// If you do not specify a provider, all errors will be retried (`true` by default)
type ShouldRetryProvider func(error) bool
