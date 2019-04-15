package logger

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	log, err := New()
	assert.Nil(err)
	assert.NotNil(log.Latch)
	assert.NotNil(log.Context)
	assert.NotNil(log.Formatter)
	assert.NotNil(log.Output)
	assert.True(log.RecoverPanics)

	for _, defaultFlag := range DefaultFlags {
		assert.True(log.Flags.IsEnabled(defaultFlag))
	}

	log, err = New(OptAll(), OptFormatter(NewJSONOutputFormatter()))
	assert.Nil(err)
	assert.True(log.Flags.IsEnabled(uuid.V4().String()))
	typed, ok := log.Formatter.(*JSONOutputFormatter)
	assert.True(ok)
	assert.NotNil(typed)
}

func TestLoggerE2ESubContext(t *testing.T) {
	assert := assert.New(t)

	output := new(bytes.Buffer)
	log, err := New(
		OptOutput(output),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	assert.Nil(err)

	scID := uuid.V4().String()
	sc := log.SubContext(scID)

	sc.Infof("this is infof")
	sc.Errorf("this is errorf")
	sc.Fatalf("this is fatalf")

	sc.Trigger(context.Background(), NewMessageEvent(Info, "this is a triggered message"))
	assert.Nil(log.Drain())

	assert.Contains(output.String(), fmt.Sprintf("[info] [%s] this is infof", scID))
	assert.Contains(output.String(), fmt.Sprintf("[error] [%s] this is errorf", scID))
	assert.Contains(output.String(), fmt.Sprintf("[fatal] [%s] this is fatalf", scID))
	assert.Contains(output.String(), fmt.Sprintf("[info] [%s] this is a triggered message", scID))
}

func TestLoggerE2ESubContextFields(t *testing.T) {
	assert := assert.New(t)

	output := new(bytes.Buffer)
	log, err := New(
		OptOutput(output),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	assert.Nil(err)

	fieldKey := uuid.V4().String()
	fieldValue := uuid.V4().String()
	sc := log.WithFields(Fields{fieldKey: fieldValue})

	sc.Infof("this is infof")
	sc.Errorf("this is errorf")
	sc.Fatalf("this is fatalf")

	sc.Trigger(context.Background(), NewMessageEvent(Info, "this is a triggered message"))
	assert.Nil(log.Drain())

	assert.Contains(output.String(), fmt.Sprintf("[info] this is infof\t%s=%s", fieldKey, fieldValue))
	assert.Contains(output.String(), fmt.Sprintf("[error] this is errorf\t%s=%s", fieldKey, fieldValue))
	assert.Contains(output.String(), fmt.Sprintf("[fatal] this is fatalf\t%s=%s", fieldKey, fieldValue))
	assert.Contains(output.String(), fmt.Sprintf("[info] this is a triggered message\t%s=%s", fieldKey, fieldValue))

}
