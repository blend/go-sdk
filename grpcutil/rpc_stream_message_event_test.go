package grpcutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestRPCStreamMessageEvent(t *testing.T) {
	assert := assert.New(t)

	re := NewRPCStreamMessageEvent("/v1.foo", StreamMessageDirectionReceive, time.Second,
		OptRPCStreamMessageAuthority("event-authority"),
		OptRPCStreamMessageContentType("event-content-type"),
		OptRPCStreamMessageElapsed(time.Millisecond),
		OptRPCStreamMessageEngine("event-engine"),
		OptRPCStreamMessageErr(fmt.Errorf("test error")),
		OptRPCStreamMessageMethod("/v1.bar"),
		OptRPCStreamMessagePeer("event-peer"),
		OptRPCStreamMessageUserAgent("event-user-agent"),
		OptRPCStreamMessageDirection(StreamMessageDirectionSend),
	)

	assert.Equal("event-authority", re.Authority)
	assert.Equal("event-content-type", re.ContentType)
	assert.Equal(time.Millisecond, re.Elapsed)
	assert.Equal("event-engine", re.Engine)
	assert.Equal(fmt.Errorf("test error"), re.Err)
	assert.Equal("/v1.bar", re.Method)
	assert.Equal("event-peer", re.Peer)
	assert.Equal("event-user-agent", re.UserAgent)
	assert.Equal(StreamMessageDirectionSend, re.Direction)

	buf := new(bytes.Buffer)
	noColor := logger.TextOutputFormatter{
		NoColor: true,
	}

	re.WriteText(noColor, buf)
	assert.Equal("[event-engine] /v1.bar send event-peer event-authority event-user-agent event-content-type 1ms failed", buf.String())

	contents, err := json.Marshal(re)
	assert.Nil(err)
	assert.Contains(string(contents), "event-engine")
}
