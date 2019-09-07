package jobkit

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/blend/go-sdk/cron"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/sh"
)

// ShellActionOption mutates a ShellActionOptions object.
type ShellActionOption func(*ShellActionOptions)

// OptShellActionDiscardOutput sets the `Discard` field on the options.
func OptShellActionDiscardOutput(discard bool) ShellActionOption {
	return func(opts *ShellActionOptions) { opts.DiscardOutput = discard }
}

// ShellActionOptions captures options for a shell action.
type ShellActionOptions struct {
	DiscardOutput bool
}

// CreateShellAction creates a new shell action.
func CreateShellAction(exec []string, opts ...ShellActionOption) func(context.Context) error {
	var options ShellActionOptions
	for _, opt := range opts {
		opt(&options)
	}

	return func(ctx context.Context) error {
		ji := cron.GetJobInvocation(ctx)
		if ji == nil || ji.Output == nil {
			return fmt.Errorf("shell action; invocation meta required with the output set")
		}
		cmd, err := sh.CmdContext(ctx, exec[0], exec[1:]...)
		if err != nil {
			return err
		}
		if !options.DiscardOutput {
			cmd.Stdout = io.MultiWriter(ji.Output, os.Stdout)
			cmd.Stderr = io.MultiWriter(ji.Output, os.Stderr)
		}
		return ex.New(cmd.Run())
	}
}