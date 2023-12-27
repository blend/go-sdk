/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stats

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func Test_NewAddListenersOptions(t *testing.T) {
	its := assert.New(t)

	its.True(NewAddListenerOptions(OptIncludeLoggerLabelsAsTags(true)).IncludeLoggerLabelsAsTags)
}

func Test_AddListenerOptions_GetLoggerLabelsAsTags(t *testing.T) {
	its := assert.New(t)

	ctx := logger.WithLabels(context.Background(), logger.Labels{
		"foo":		"bar",
		"not-foo":	"not-bar",
	})

	its.Empty(AddListenerOptions{}.GetLoggerLabelsAsTags(ctx))
	its.NotEmpty(AddListenerOptions{
		IncludeLoggerLabelsAsTags: true,
	}.GetLoggerLabelsAsTags(ctx))
}
