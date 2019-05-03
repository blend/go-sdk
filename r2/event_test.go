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
	e.Started = time.Date(2019, 05, 02, 12, 13, 14, 15, time.UTC)
	e.Timestamp = time.Date(2019, 05, 02, 12, 13, 15, 15, time.UTC)
	output2 := new(bytes.Buffer)
	e.WriteText(logger.NewTextOutputFormatter(logger.OptTextNoColor()), output2)
	assert.Equal("GET http://localhost/http://test.com 200 (1s)\nfoo", output2.String())
}

func TestEventMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	e := NewEvent(Flag,
		OptEventStarted(time.Date(2019, 05, 02, 12, 13, 14, 15, time.UTC)),
		OptEventCompleted(time.Date(2019, 05, 02, 12, 13, 15, 15, time.UTC)),
		OptEventRequest(webutil.NewMockRequest("GET", "/foo")),
		OptEventResponse(&http.Response{StatusCode: http.StatusOK, ContentLength: 500}),
		OptEventBody([]byte("foo")),
	)

	contents, err := e.MarshalJSON()
	assert.Nil(err)
	assert.NotEmpty(contents)

	var jsonContents struct {
		Req struct {
			startTime time.Time
			Method    string `json:"method"`
			URL       string `json:"url"`
			headers   map[string][]string
		} `json:"req"`
		Res struct {
			completeTime  time.Time
			StatusCode    int `json:"statusCode"`
			ContentLength int `json:"contentLength"`
			headers       map[string][]string
		} `json:"res"`
	}

	assert.Nil(json.Unmarshal(contents, &jsonContents))
	assert.Equal("http://localhost/foo", jsonContents.Req.URL)
	assert.Equal("GET", jsonContents.Req.Method)
	assert.Equal(http.StatusOK, jsonContents.Res.StatusCode)
	assert.Equal(500, jsonContents.Res.ContentLength)
}
