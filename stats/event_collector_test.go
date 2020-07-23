package stats

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestEventCheck(t *testing.T) {
	assert := assert.New(t)

	event := Event{}
	err := event.Check()
	assert.NotNil(event.Check())
	assert.Contains(err.Error(), "event title")

	event.Text = "text"
	err = event.Check()
	assert.NotNil(event.Check())
	assert.Contains(err.Error(), "event title")

	event.Title = "title"
	event.Text = ""
	err = event.Check()
	assert.NotNil(event.Check())
	assert.Contains(err.Error(), "event text")

	event.Text = "text"
	err = event.Check()
	assert.Nil(err)
	assert.Nil(event.Check())
}
