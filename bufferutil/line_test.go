package bufferutil

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestFilterLines(t *testing.T) {
	assert := assert.New(t)

	lines := []Line{
		{Timestamp: time.Date(2019, 9, 4, 12, 11, 10, 9, time.UTC), Line: []byte("this is just a test")},
		{Timestamp: time.Date(2019, 9, 5, 12, 11, 10, 9, time.UTC), Line: []byte("this is just another test")},
	}
	filtered := FilterLines(lines, func(l Line) bool {
		return l.Timestamp.After(lines[0].Timestamp)
	})
	assert.Len(filtered, 1)
}

func TestLineBufferLinesJSON(t *testing.T) {
	assert := assert.New(t)

	lines := []Line{
		{Timestamp: time.Date(2019, 9, 4, 12, 11, 10, 9, time.UTC), Line: []byte("this is just a test")},
		{Timestamp: time.Date(2019, 9, 5, 12, 11, 10, 9, time.UTC), Line: []byte("this is just another test")},
	}

	contents, err := json.Marshal(lines)
	assert.Nil(err)

	var verify []Line
	err = json.Unmarshal(contents, &verify)
	assert.Nil(err)
	assert.Len(verify, 2)
}
