package jsonutil_test

import (
	"reflect"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/jsonutil"
)

func TestGet(t *testing.T) {
	it := assert.New(t)

	values := map[string]interface{}{
		"hey": "hello",
		"bye": map[string]interface{}{
			"seeya": 18,
		},
	}
	// Empty path.
	value, exists := jsonutil.Get(values)
	it.Nil(value)
	it.False(exists)
	// Non-interface
	value, exists = jsonutil.Get(48, "once", "twice")
	it.Nil(value)
	it.False(exists)
	// Simple path.
	value, exists = jsonutil.Get(values, "hey")
	it.Equal("hello", value)
	it.True(exists)
	// Nested path.
	value, exists = jsonutil.Get(values, "bye", "seeya")
	it.Equal(18, value)
	it.True(exists)
}

func TestGetMap(t *testing.T) {
	it := assert.New(t)

	values := map[string]interface{}{
		"over": "here",
		"under": map[string]interface{}{
			"there": false,
		},
	}

	// Empty path.
	value, exists := jsonutil.GetMap(values)
	it.Nil(value)
	it.False(exists)

	// Missing nested path.
	value, exists = jsonutil.GetMap(values, "nope", "nope")
	it.Nil(value)
	it.False(exists)

	// FirstKey, single path.
	value, exists = jsonutil.GetMap(values, jsonutil.FirstKey)
	// NOTE: `FirstKey` is non-deterministic since it relies on `map` order.
	it.True(reflect.DeepEqual("here", value) || reflect.DeepEqual(values["under"], value))
	it.True(exists)

	// FirstKey, empty map.
	value, exists = jsonutil.GetMap(map[string]interface{}{}, jsonutil.FirstKey)
	it.Nil(value)
	it.False(exists)

	// LastKey, single path.
	value, exists = jsonutil.GetMap(values, jsonutil.LastKey)
	// NOTE: `LastKey` is non-deterministic since it relies on `map` order.
	it.True(reflect.DeepEqual("here", value) || reflect.DeepEqual(values["under"], value))
	it.True(exists)

	// LastKey, empty map.
	value, exists = jsonutil.GetMap(map[string]interface{}{}, jsonutil.LastKey)
	it.Nil(value)
	it.False(exists)
}
