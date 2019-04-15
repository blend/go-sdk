package cron

import "github.com/blend/go-sdk/ex"

const (
	// ErrJobNotLoaded is a common error.
	ErrJobNotLoaded ex.Class = "job not loaded"

	// ErrJobAlreadyLoaded is a common error.
	ErrJobAlreadyLoaded ex.Class = "job already loaded"

	// ErrJobNotFound is a common error.
	ErrJobNotFound ex.Class = "job not found"

	// ErrJobCancelled is a common error.
	ErrJobCancelled ex.Class = "job cancelled"
)

// IsJobNotLoaded returns if the error is a job not loaded error.
func IsJobNotLoaded(err error) bool {
	return ex.Is(err, ErrJobNotLoaded)
}

// IsJobAlreadyLoaded returns if the error is a job already loaded error.
func IsJobAlreadyLoaded(err error) bool {
	return ex.Is(err, ErrJobAlreadyLoaded)
}

// IsJobNotFound returns if the error is a task not found error.
func IsJobNotFound(err error) bool {
	return ex.Is(err, ErrJobNotFound)
}

// IsJobCancelled returns if the error is a task not found error.
func IsJobCancelled(err error) bool {
	return ex.Is(err, ErrJobCancelled)
}
