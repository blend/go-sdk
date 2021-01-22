package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/blend/go-sdk/logger"
)

// FailureCodes
const (
	SuiteFailureTests  = 1
	SuiteFailureBefore = 2
	SuiteFailureAfter  = 3
)

// New returns a new test suite.
func New(m *testing.M, opts ...Option) *Suite {
	s := Suite{
		M: m,
	}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

// Option is a mutator for a test suite.
type Option func(*Suite)

// SuiteAction is a step that can be run either before or after package tests.
type SuiteAction func(context.Context) error

// Suite is a set of before and after actions for a given package tests.
type Suite struct {
	M      *testing.M
	Log    logger.Log
	Before []SuiteAction
	After  []SuiteAction
}

// Run runs tests and returns the exit code.
func (s Suite) Run() {
	var code int
	defer func() {
		os.Exit(code)
	}()
	ctx := context.Background()
	if s.Log != nil {
		ctx = logger.WithLogger(ctx, s.Log)
	}
	var err error
	for _, before := range s.Before {
		if err = before(ctx); err != nil {
			logger.MaybeFatalf(s.Log, "error during setup steps: %+v", err)
			code = SuiteFailureBefore
			return
		}
	}
	defer func() {
		for _, after := range s.After {
			if err = after(ctx); err != nil {
				logger.MaybeFatalf(s.Log, "error during cleanup steps: %+v", err)
				code = SuiteFailureAfter
				return
			}
		}
	}()
	code = s.M.Run()
}
