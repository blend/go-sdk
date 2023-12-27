/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// FailureCodes
const (
	SuiteFailureTests	= 1
	SuiteFailureBefore	= 2
	SuiteFailureAfter	= 3
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
	M	*testing.M
	Log	logger.Log
	Before	[]SuiteAction
	After	[]SuiteAction
}

// Run runs tests and calls os.Exit(...) with the exit code.
func (s Suite) Run() {
	os.Exit(s.RunCode())
}

// RunCode runs the suite and returns an exit code.
//
// It is used by `.Run()`, which will os.Exit(...) this code.
func (s Suite) RunCode() (code int) {
	ctx := context.Background()
	if s.Log != nil {
		ctx = logger.WithLogger(ctx, s.Log)
	}
	var err error
	for _, before := range s.Before {
		if err = executeSafe(ctx, before); err != nil {
			logger.MaybeFatalf(s.Log, "error during setup steps: %+v", err)
			code = SuiteFailureBefore
			return
		}
	}
	defer func() {
		for _, after := range s.After {
			if err = executeSafe(ctx, after); err != nil {
				logger.MaybeFatalf(s.Log, "error during cleanup steps: %+v", err)
				code = SuiteFailureAfter
				return
			}
		}
	}()
	if s.M != nil {
		code = s.M.Run()
	}
	return
}

func executeSafe(ctx context.Context, action func(context.Context) error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ex.New(r)
		}
	}()
	err = action(ctx)
	return
}
