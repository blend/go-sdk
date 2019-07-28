package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func value(value string, ok bool) string {
	return value
}

func TestAnnotations(t *testing.T) {
	assert := assert.New(t)

	a := make(Annotations)
	assert.Empty(a)
	a.AddAnnotationValue("foo", "bar")
	assert.NotEmpty(a)
	assert.Equal("bar", value(a.GetAnnotationValue("foo")))
	assert.Empty(a.GetAnnotationValue("bar"))

	a.AddAnnotationValue("buzz", "fuzz")
	assert.Equal("fuzz", value(a.GetAnnotationValue("buzz")))
	assert.Equal("bar", value(a.GetAnnotationValue("foo")))

	assert.Any(a.GetAnnotationKeys(), func(v interface{}) bool {
		return v.(string) == "foo"
	})
	assert.Any(a.GetAnnotationKeys(), func(v interface{}) bool {
		return v.(string) == "buzz"
	})

	values := a.Decompose()
	assert.NotEmpty(values)

	assert.Equal("bar", values["foo"])
	assert.Equal("fuzz", values["buzz"])
}
