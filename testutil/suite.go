/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package testutil

import (
	"context"
	"fmt"
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
		code = SuiteFailureBefore
		defer recoverHookPanicAndExitWithCode(code)
		if err = before(ctx); err != nil {
			logger.MaybeFatalf(s.Log, "error during setup steps: %+v", err)
			return
		}
	}
	defer func() {
		for _, after := range s.After {
			code = SuiteFailureAfter
			defer recoverHookPanicAndExitWithCode(code)
			if err = after(ctx); err != nil {
				logger.MaybeFatalf(s.Log, "error during cleanup steps: %+v", err)
				return
			}
		}
	}()
	code = s.M.Run()
}

// recovers and exits with the given code.
// meant specifically for panics in an OptBefore or OptAfter hook
func recoverHookPanicAndExitWithCode(code int) {
	var hookPhase string
	if code == SuiteFailureBefore {
		hookPhase = "before"
	}
	if code == SuiteFailureAfter {
		hookPhase = "after"
	}
	if r := recover(); r != nil {
		fmt.Printf("a panic occured in one of the %s hooks: %+v\n", hookPhase, r)
		os.Exit(code)
	}
}
