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
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestQueryEvent(t *testing.T) {
	assert := assert.New(t)

	qe := NewQueryEvent("query-body", time.Second,
		OptQueryEventBody("event-body"),
		OptQueryEventDatabase("event-database"),
		OptQueryEventEngine("event-engine"),
		OptQueryEventUsername("event-username"),
		OptQueryEventLabel("event-query-label"),
		OptQueryEventElapsed(time.Millisecond),
		OptQueryEventErr(fmt.Errorf("test error")),
	)

	assert.Equal("event-body", qe.Body)
	assert.Equal("event-database", qe.Database)
	assert.Equal("event-engine", qe.Engine)
	assert.Equal("event-username", qe.Username)
	assert.Equal("event-query-label", qe.Label)
	assert.Equal(time.Millisecond, qe.Elapsed)
	assert.Equal("test error", qe.Err.Error())

	buf := new(bytes.Buffer)
	noColor := logger.TextOutputFormatter{
		NoColor: true,
	}

	qe.WriteText(noColor, buf)
	assert.Equal("[event-engine event-username@event-database] [event-query-label] event-body 1ms failed", buf.String())

	contents, err := json.Marshal(qe)
	assert.Nil(err)
	assert.Contains(string(contents), "event-engine")
}

func TestQueryEventListener(t *testing.T) {
	assert := assert.New(t)

	qe := NewQueryEvent("select 1", time.Second)

	var didCall bool
	ml := NewQueryEventListener(func(ctx context.Context, ae QueryEvent) {
		didCall = true
	})
	ml(context.Background(), qe)
	assert.True(didCall)
}
