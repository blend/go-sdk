package r2

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

func TestNewEvent(t *testing.T) {
	assert := assert.New(t)

	e := NewEvent(Flag, OptEventBody([]byte("foo")))
	assert.Equal("foo", e.Body)
}

func TestEventWriteString(t *testing.T) {
	assert := assert.New(t)

	e := NewEvent(Flag,
		OptEventRequest(webutil.NewMockRequest("GET", "http://test.com")),
		OptEventBody([]byte("foo")),
	)

	output := new(bytes.Buffer)
	e.WriteText(logger.NewTextOutputFormatter(logger.OptTextNoColor()), output)
	assert.Equal("GET http://localhost/http://test.com\nfoo", output.String())

	e.Response = &http.Response{
		StatusCode: http.StatusOK,
	}
	e.Elapsed = time.Second
	output2 := new(bytes.Buffer)
	e.WriteText(logger.NewTextOutputFormatter(logger.OptTextNoColor()), output2)
	assert.Equal("GET http://localhost/http://test.com 200 (1s)\nfoo", output2.String())
}

// eventJSONSchema is the json schema of the logger event.
type eventJSONSchema struct {
	Req struct {
		StartTime time.Time           `json:"startTime"`
		Method    string              `json:"method"`
		URL       string              `json:"url"`
		Headers   map[string][]string `json:"headers"`
	} `json:"req"`
	Res struct {
		CompleteTime  time.Time           `json:"completeTime"`
		StatusCode    int                 `json:"statusCode"`
		ContentLength int                 `json:"contentLength"`
		Headers       map[string][]string `json:"headers"`
	} `json:"res"`
	Body string `json:"body"`
}

func TestEventMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	e := NewEvent(Flag,
		OptEventRequest(webutil.NewMockRequest("GET", "/foo")),
		OptEventResponse(&http.Response{StatusCode: http.StatusOK, ContentLength: 500}),
		OptEventBody([]byte("foo")),
	)

	contents, err := json.Marshal(e.Decompose())
	assert.Nil(err)
	assert.NotEmpty(contents)

	var jsonContents eventJSONSchema
	assert.Nil(json.Unmarshal(contents, &jsonContents))
	assert.Equal("http://localhost/foo", jsonContents.Req.URL)
	assert.Equal("GET", jsonContents.Req.Method)
	assert.Equal(http.StatusOK, jsonContents.Res.StatusCode)
	assert.Equal(500, jsonContents.Res.ContentLength)
	assert.Equal("foo", jsonContents.Body)
}
