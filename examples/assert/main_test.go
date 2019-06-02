package main_test

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestExample(t *testing.T) {
	assert := assert.New(t)

	assert.True(true)
	assert.False(false)
	assert.Equal("foo", "foo")
	assert.NotEqual("foo", "bar")
	assert.Any([]int{1, 2, 3}, func(v interface{}) bool { return v.(int) == 1 })
	assert.All([]int{1, 2, 3}, func(v interface{}) bool { return v.(int) > 0 })
	assert.None([]int{1, 2, 3}, func(v interface{}) bool { return v.(int) == 0 })
}

func TestExampleFailure(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("foo", "bar")
}

func TestFilteredUnit(t *testing.T) {
	assert := assert.Filtered(t, assert.Unit)

	assert.True(true)
	assert.False(false)
	assert.Equal("foo", "foo")
	assert.NotEqual("foo", "bar")
	assert.Any([]int{1, 2, 3}, func(v interface{}) bool { return v.(int) == 1 })
	assert.All([]int{1, 2, 3}, func(v interface{}) bool { return v.(int) > 0 })
	assert.None([]int{1, 2, 3}, func(v interface{}) bool { return v.(int) == 0 })
}

func TestFilteredAcceptance(t *testing.T) {
	assert := assert.Filtered(t, assert.Acceptance)

	assert.True(true)
	assert.False(false)
	assert.Equal("foo", "foo")
	assert.NotEqual("foo", "bar")
	assert.Any([]int{1, 2, 3}, func(v interface{}) bool { return v.(int) == 1 })
	assert.All([]int{1, 2, 3}, func(v interface{}) bool { return v.(int) > 0 })
	assert.None([]int{1, 2, 3}, func(v interface{}) bool { return v.(int) == 0 })
}

func TestFilteredIntegration(t *testing.T) {
	assert := assert.Filtered(t, assert.Integration)

	assert.True(true)
	assert.False(false)
	assert.Equal("foo", "foo")
	assert.NotEqual("foo", "bar")
	assert.Any([]int{1, 2, 3}, func(v interface{}) bool { return v.(int) == 1 })
	assert.All([]int{1, 2, 3}, func(v interface{}) bool { return v.(int) > 0 })
	assert.None([]int{1, 2, 3}, func(v interface{}) bool { return v.(int) == 0 })
}
