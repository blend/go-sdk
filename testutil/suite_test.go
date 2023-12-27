/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package testutil

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Suite_panics_before(t *testing.T) {
	its := assert.New(t)

	panicsBefore := Suite{
		Before: []SuiteAction{
			func(_ context.Context) error { return nil },
			func(_ context.Context) error { panic("at the disco") },
		},
		After: []SuiteAction{
			func(_ context.Context) error { return nil },
		},
	}
	its.Equal(SuiteFailureBefore, panicsBefore.RunCode())
}

func Test_Suite_panics_after(t *testing.T) {
	its := assert.New(t)

	panicsBefore := Suite{
		Before: []SuiteAction{
			func(_ context.Context) error { return nil },
		},
		After: []SuiteAction{
			func(_ context.Context) error { return nil },
			func(_ context.Context) error { panic("at the disco") },
		},
	}
	its.Equal(SuiteFailureAfter, panicsBefore.RunCode())
}
