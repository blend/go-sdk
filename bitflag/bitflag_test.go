package bitflag

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCombine(t *testing.T) {
	assert := assert.New(t)
	three := Combine(1, 2)
	assert.Equal(3, three)
}

func TestAny(t *testing.T) {
	assert := assert.New(t)

	var one Bitflag = 1 << 0
	var two Bitflag = 1 << 1
	var four Bitflag = 1 << 2
	var eight Bitflag = 1 << 3
	var sixteen Bitflag = 1 << 4

	masterFlag := Combine(one, two, four, eight)
	checkFlag := Combine(one, sixteen)
	assert.True(masterFlag.Any(checkFlag))
	assert.False(masterFlag.Any(1 << 5))
}

func TestBitFlagAll(t *testing.T) {
	assert := assert.New(t)

	var one Bitflag = 1 << 0
	var two Bitflag = 1 << 1
	var four Bitflag = 1 << 2
	var eight Bitflag = 1 << 3
	var sixteen Bitflag = 1 << 4

	masterFlag := Combine(one, two, four, eight)
	checkValidFlag := Combine(one, two)
	checkInvalidFlag := Combine(one, sixteen)
	assert.True(masterFlag.All(checkValidFlag))
	assert.False(masterFlag.All(checkInvalidFlag))
}

func TestBitFlagSet(t *testing.T) {
	assert := assert.New(t)

	var zero Bitflag
	var four Bitflag = 1 << 4
	assert.Equal(16, zero.Set(1<<4))
	assert.Equal(16, four.Set(1<<4))
}

func TestBitFlagZero(t *testing.T) {
	assert := assert.New(t)

	var zero Bitflag
	var one Bitflag = 1 << 0
	var four Bitflag = 1 << 4

	flagSet := Combine(one, four)
	assert.Equal(0, zero.Unset(1<<4))
	assert.Equal(0, four.Unset(1<<4))
	assert.Equal(16, flagSet.Unset(one))
}
