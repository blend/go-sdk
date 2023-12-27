/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	stringSource := String("")

	boolValue, err := Parse(stringSource).Bool(context.TODO())
	assert.Nil(err)
	assert.Nil(boolValue)

	trueValues := []string{"1", "true", "yes", "on"}
	for _, tv := range trueValues {
		stringSource = String(tv)
		boolValue, err = Parse(stringSource).Bool(context.TODO())
		assert.Nil(err)
		assert.NotNil(boolValue)
		assert.True(*boolValue)
	}

	falseValues := []string{"0", "false", "no", "off"}
	for _, fv := range falseValues {
		stringSource = String(fv)
		boolValue, err = Parse(stringSource).Bool(context.TODO())
		assert.Nil(err)
		assert.NotNil(boolValue)
		assert.False(*boolValue)
	}

	stringSource = String("not a bool")
	boolValue, err = Parse(stringSource).Bool(context.TODO())
	assert.NotNil(err)
	assert.Nil(boolValue)

	stringSource = String("")
	intValue, err := Parse(stringSource).Int(context.TODO())
	assert.Nil(err)
	assert.Nil(intValue)

	stringSource = String("bad value")
	intValue, err = Parse(stringSource).Int(context.TODO())
	assert.NotNil(err)
	assert.Nil(intValue)

	stringSource = String("1234")
	intValue, err = Parse(stringSource).Int(context.TODO())
	assert.Nil(err)
	assert.NotNil(intValue)
	assert.Equal(1234, *intValue)

	stringSource = String("")
	floatValue, err := Parse(stringSource).Float64(context.TODO())
	assert.Nil(err)
	assert.Nil(floatValue)

	stringSource = String("bad value")
	floatValue, err = Parse(stringSource).Float64(context.TODO())
	assert.NotNil(err)
	assert.Nil(floatValue)

	stringSource = String("1234.34")
	floatValue, err = Parse(stringSource).Float64(context.TODO())
	assert.Nil(err)
	assert.NotNil(floatValue)
	assert.Equal(1234.34, *floatValue)

	stringSource = String("")
	durationValue, err := Parse(stringSource).Duration(context.TODO())
	assert.Nil(err)
	assert.Nil(durationValue)

	stringSource = String("bad value")
	durationValue, err = Parse(stringSource).Duration(context.TODO())
	assert.NotNil(err)
	assert.Nil(durationValue)

	stringSource = String("10s")
	durationValue, err = Parse(stringSource).Duration(context.TODO())
	assert.Nil(err)
	assert.NotNil(durationValue)
	assert.Equal(10*time.Second, *durationValue)
}
