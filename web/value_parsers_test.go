package web

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_BoolValue(t *testing.T) {
	its := assert.New(t)

	var value bool
	var err error

	testErr := fmt.Errorf("test error")
	value, err = BoolValue("", testErr)
	its.Equal(testErr, err)
	its.False(value)

	trueValues := []string{"1", "true", "yes", "on"}
	for _, tv := range trueValues {
		value, err = BoolValue(tv, nil)
		its.Nil(err)
		its.True(value)
	}

	falseValues := []string{"0", "false", "no", "off"}
	for _, tv := range falseValues {
		value, err = BoolValue(tv, nil)
		its.Nil(err)
		its.False(value)
	}

	value, err = BoolValue("garbage", nil)
	its.Equal(ErrInvalidBoolValue, err)
	its.False(value)
}

func Test_IntValue(t *testing.T) {
	its := assert.New(t)

	var value int
	var err error

	testErr := fmt.Errorf("test error")
	value, err = IntValue("", testErr)
	its.Equal(testErr, err)
	its.Zero(value)

	value, err = IntValue("1234", nil)
	its.Nil(err)
	its.Equal(1234, value)

	value, err = IntValue("garbage", nil)
	its.NotNil(err)
	its.Zero(value)
}

func Test_Int64Value(t *testing.T) {
	its := assert.New(t)

	var value int64
	var err error

	testErr := fmt.Errorf("test error")
	value, err = Int64Value("", testErr)
	its.Equal(testErr, err)
	its.Zero(value)

	value, err = Int64Value("1234", nil)
	its.Nil(err)
	its.Equal(1234, value)

	value, err = Int64Value("garbage", nil)
	its.NotNil(err)
	its.Zero(value)
}

func Test_Float64Value(t *testing.T) {
	its := assert.New(t)

	var value float64
	var err error

	testErr := fmt.Errorf("test error")
	value, err = Float64Value("", testErr)
	its.Equal(testErr, err)
	its.Zero(value)

	value, err = Float64Value("1234.56", nil)
	its.Nil(err)
	its.Equal(1234.56, value)

	value, err = Float64Value("garbage", nil)
	its.NotNil(err)
	its.Zero(value)
}

func Test_DurationValue(t *testing.T) {
	its := assert.New(t)

	var value time.Duration
	var err error

	testErr := fmt.Errorf("test error")
	value, err = DurationValue("", testErr)
	its.Equal(testErr, err)
	its.Zero(value)

	value, err = DurationValue("10s", nil)
	its.Nil(err)
	its.Equal(10*time.Second, value)

	value, err = DurationValue("garbage", nil)
	its.NotNil(err)
	its.Zero(value)
}

func Test_StringValue(t *testing.T) {
	its := assert.New(t)

	var value string
	testErr := fmt.Errorf("test error")
	value = StringValue("foo", testErr)
	its.Equal("foo", value)
}

func Test_CSVValue(t *testing.T) {
	its := assert.New(t)

	var value []string
	var err error

	testErr := fmt.Errorf("test error")
	value, err = CSVValue("", testErr)
	its.Equal(testErr, err)
	its.Empty(value)

	value, err = CSVValue("foo,bar", nil)
	its.Nil(err)
	its.Equal([]string{"foo", "bar"}, value)
}
