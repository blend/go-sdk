package logger

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/ansi"

	"github.com/blend/go-sdk/assert"
)

func TestNewEventMeta(t *testing.T) {
	assert := assert.New(t)

	em := NewEventMeta(Info)
	assert.Equal(Info, em.GetFlag())
	assert.False(em.GetTimestamp().IsZero())
}

func TestEventMetaOptions(t *testing.T) {
	assert := assert.New(t)

	em := NewEventMeta(Info)
	assert.Equal(Info, em.GetFlag())
	OptEventMetaFlag(Error)(em)
	assert.Equal(Error, em.GetFlag())

	assert.False(em.GetTimestamp().IsZero())
	OptEventMetaTimestamp(time.Time{})(em)
	assert.True(em.GetTimestamp().IsZero())

	assert.Empty(em.GetFlagColor())
	OptEventMetaFlagColor(ansi.ColorBlue)(em)
	assert.Equal(ansi.ColorBlue, em.GetFlagColor())
}

func TestEventMetaDecompose(t *testing.T) {
	assert := assert.New(t)

	decomposed := NewEventMeta(Info).Decompose()
	assert.Equal("info", decomposed[FieldFlag])
}
