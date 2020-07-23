package envoyutil_test

import (
	"encoding/json"
	"testing"

	sdkAssert "github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"

	"github.com/blend/go-sdk/envoyutil"
)

func TestXFCCExtractionErrorMarshal(t *testing.T) {
	assert := sdkAssert.New(t)

	c := ex.Class("caused by something invalid")
	err := &envoyutil.XFCCExtractionError{Class: c, XFCC: "a=b", Metadata: map[string]string{"x": "why"}}

	asBytes, marshalErr := json.MarshalIndent(err, "", "  ")
	assert.Nil(marshalErr)
	expected := `{
  "class": "caused by something invalid",
  "xfcc": "a=b",
  "metadata": {
    "x": "why"
  }
}`
	assert.Equal(expected, string(asBytes))
}

func TestXFCCExtractionErrorError(t *testing.T) {
	assert := sdkAssert.New(t)

	c := ex.Class("oh a bad thing happened")
	var err error = &envoyutil.XFCCExtractionError{Class: c}
	assert.Equal(c, err.Error())
}

func TestXFCCFatalErrorMarshal(t *testing.T) {
	assert := sdkAssert.New(t)

	c := ex.Class("caused by something fatal")
	err := &envoyutil.XFCCFatalError{Class: c, XFCC: "c=d"}

	asBytes, marshalErr := json.MarshalIndent(err, "", "  ")
	assert.Nil(marshalErr)
	expected := `{
  "class": "caused by something fatal",
  "xfcc": "c=d"
}`
	assert.Equal(expected, string(asBytes))
}

func TestXFCCFatalErrorError(t *testing.T) {
	assert := sdkAssert.New(t)

	c := ex.Class("oh a fatal thing happened")
	var err error = &envoyutil.XFCCFatalError{Class: c}
	assert.Equal(c, err.Error())
}
