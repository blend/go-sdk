package statsd

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_Metric_Float64(t *testing.T) {
	assert := assert.New(t)

	good := Metric{Value: "3.14"}
	goodValue, err := good.Float64()
	assert.Nil(err)
	assert.Equal(3.14, goodValue)

	bad := Metric{Value: "foo"}
	badValue, err := bad.Float64()
	assert.NotNil(err)
	assert.Zero(badValue)
}

func Test_Metric_Int64(t *testing.T) {
	assert := assert.New(t)

	good := Metric{Value: "314"}
	goodValue, err := good.Int64()
	assert.Nil(err)
	assert.Equal(314, goodValue)

	bad := Metric{Value: "foo"}
	badValue, err := bad.Int64()
	assert.NotNil(err)
	assert.Zero(badValue)
}

func Test_Metric_Duration(t *testing.T) {
	assert := assert.New(t)

	good := Metric{Value: "512.12"}
	goodValue, err := good.Duration()
	assert.Nil(err)
	assert.Equal(time.Duration(512.12*float64(time.Millisecond)), goodValue)

	bad := Metric{Value: "foo"}
	badValue, err := bad.Duration()
	assert.NotNil(err)
	assert.Zero(badValue)
}
