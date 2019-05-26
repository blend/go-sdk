package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestQueryEvent(t *testing.T) {
	assert := assert.New(t)

	qe := NewQueryEvent("query-body", time.Second,
		OptQueryBody("event-body"),
		OptQueryDatabase("event-database"),
		OptQueryEngine("event-engine"),
		OptQueryUsername("event-username"),
		OptQueryLabel("event-query-label"),
		OptQueryElapsed(time.Millisecond),
		OptQueryErr(fmt.Errorf("test error")),
	)

	assert.Equal("event-body", qe.Body)
	assert.Equal("event-database", qe.Database)
	assert.Equal("event-engine", qe.Engine)
	assert.Equal("event-username", qe.Username)
	assert.Equal("event-query-label", qe.QueryLabel)
	assert.Equal(time.Millisecond, qe.Elapsed)
	assert.Equal("test error", qe.Err.Error())

	buf := new(bytes.Buffer)
	noColor := TextOutputFormatter{
		NoColor: true,
	}

	qe.WriteText(noColor, buf)
	assert.Equal("[event-engine event-username@event-database] [event-query-label] 1ms failed event-body", buf.String())

	contents, err := json.Marshal(qe)
	assert.Nil(err)
	assert.Contains(string(contents), "event-engine")
}
