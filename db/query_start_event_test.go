/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func Test_NewQueryStartEvent(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	qe := NewQueryStartEvent("query-body",
		OptQueryStartEventBody("event-body"),
		OptQueryStartEventDatabase("event-database"),
		OptQueryStartEventEngine("event-engine"),
		OptQueryStartEventUsername("event-username"),
		OptQueryStartEventLabel("event-query-label"),
	)

	its.Equal("event-body", qe.Body)
	its.Equal("event-database", qe.Database)
	its.Equal("event-engine", qe.Engine)
	its.Equal("event-username", qe.Username)
	its.Equal("event-query-label", qe.Label)

	buf := new(bytes.Buffer)
	noColor := logger.TextOutputFormatter{
		NoColor: true,
	}

	qe.WriteText(noColor, buf)
	its.Equal("[event-engine event-username@event-database] [event-query-label] event-body", buf.String())

	contents, err := json.Marshal(qe)
	its.Nil(err)
	its.Contains(string(contents), "event-engine")
}

func Test_QueryStartEventListener(t *testing.T) {
	assert := assert.New(t)

	qe := NewQueryStartEvent("select 1")

	var didCall bool
	ml := NewQueryStartEventListener(func(ctx context.Context, ae QueryStartEvent) {
		didCall = true
	})
	ml(context.Background(), qe)
	assert.True(didCall)
}
