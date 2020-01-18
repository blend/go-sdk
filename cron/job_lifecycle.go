package cron

import "context"

// JobLifecycle is a suite of lifeycle hooks
// you can set for a given job.
type JobLifecycle struct {
	OnLoad   func() error
	OnUnload func() error

	OnBegin        func(context.Context)
	OnCancellation func(context.Context)
	OnError        func(context.Context)
	OnComplete     func(context.Context)
	OnBroken       func(context.Context)
	OnFixed        func(context.Context)
	OnEnabled      func(context.Context)
	OnDisabled     func(context.Context)

	RestoreHistory func(context.Context) ([]JobInvocation, error)
	PersistHistory func(context.Context, []JobInvocation) error
}
