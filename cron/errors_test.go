/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func Test_Errors(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	its.True(IsJobNotLoaded(ex.New(ErrJobNotLoaded)))
	its.False(IsJobNotLoaded(ex.New("incorrect")))

	its.True(IsJobAlreadyLoaded(ex.New(ErrJobAlreadyLoaded)))
	its.False(IsJobAlreadyLoaded(ex.New("incorrect")))

	its.True(IsJobNotFound(ex.New(ErrJobNotFound)))
	its.False(IsJobNotFound(ex.New("incorrect")))

	its.True(IsJobCanceled(ex.New(ErrJobCanceled)))
	its.False(IsJobCanceled(ex.New("incorrect")))

	its.True(IsJobAlreadyRunning(ex.New(ErrJobAlreadyRunning)))
	its.False(IsJobAlreadyRunning(ex.New("incorrect")))
}
