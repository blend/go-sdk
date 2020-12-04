package datadog

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/ref"

	"github.com/DataDog/datadog-go/statsd"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/uuid"
)

func TestConvertEvent(t *testing.T) {
	assert := assert.New(t)

	original := stats.Event{
		Title:          uuid.V4().String(),
		Text:           uuid.V4().String(),
		Timestamp:      time.Now().UTC(),
		Hostname:       uuid.V4().String(),
		AggregationKey: uuid.V4().String(),
		Priority:       uuid.V4().String(),
		SourceTypeName: uuid.V4().String(),
		AlertType:      uuid.V4().String(),
		Tags:           []string{uuid.V4().String()},
	}

	converted := ConvertEvent(original)
	assert.Equal(original.Title, converted.Title)
	assert.Equal(original.Text, converted.Text)
	assert.Equal(original.Timestamp, converted.Timestamp)
	assert.Equal(original.Hostname, converted.Hostname)
	assert.Equal(original.AggregationKey, converted.AggregationKey)
	assert.Equal(original.Priority, converted.Priority)
	assert.Equal(original.SourceTypeName, converted.SourceTypeName)
	assert.Equal(original.AlertType, converted.AlertType)
	assert.Equal(original.Tags, converted.Tags)
}

func TestCollectorFlush(t *testing.T) {
	assert := assert.New(t)

	// `client` is `nil`
	c := Collector{}
	assert.Nil(c.Flush())

	// `client` is not `nil`
	client, err := statsd.New("localhost:8125")
	assert.Nil(err)
	defer client.Close()

	c = Collector{client: client}
	assert.Nil(c.Flush())
}

func TestCollectorClose(t *testing.T) {
	assert := assert.New(t)

	// `client` is `nil`
	c := Collector{}
	assert.Nil(c.Close())

	// `client` is not `nil`
	client, err := statsd.New("localhost:8125")
	assert.Nil(err)

	c = Collector{client: client}
	assert.Nil(c.Close())
}

func TestCollectorNew(t *testing.T) {
	assert := assert.New(t)
	cfg := Config{
		Address: "localhost:8125",
	}
	c, err := New(cfg)
	assert.Nil(err)
	assert.NotNil(c)
	assert.NotNil(c.client)
	assert.NotEmpty(c.defaultTags)

	c, err = New(cfg, statsd.WithNamespace("hello"))
	assert.Nil(err)
	assert.NotNil(c)
	assert.NotNil(c.client)
	assert.Equal("hello.", c.client.Namespace)

	cfg.BufferDepth = 10
	cfg.Buffered = ref.Bool(true)

	c, err = New(cfg, statsd.WithNamespace("hello"))
	assert.Nil(err)
	assert.NotNil(c)
	assert.NotNil(c.client)
	assert.Equal("hello.", c.client.Namespace)
}
