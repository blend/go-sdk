/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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

	c := ex.Class("caused by bad extraction")
	err := &envoyutil.XFCCExtractionError{Class: c, XFCC: "a=b", Metadata: map[string]string{"x": "why"}}

	asBytes, marshalErr := json.MarshalIndent(err, "", "  ")
	assert.Nil(marshalErr)
	expected := `{
  "class": "caused by bad extraction",
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

func TestIsExtractionError(t *testing.T) {
	assert := sdkAssert.New(t)

	var err error = ex.New("NOPE")
	assert.False(envoyutil.IsExtractionError(err))
	err = &envoyutil.XFCCExtractionError{Class: "YEP"}
	assert.True(envoyutil.IsExtractionError(err))
}

func TestXFCCValidationErrorMarshal(t *testing.T) {
	assert := sdkAssert.New(t)

	c := ex.Class("caused by something invalid")
	err := &envoyutil.XFCCValidationError{Class: c, XFCC: "mm=hm", Metadata: map[string]string{"ecks": "y"}}

	asBytes, marshalErr := json.MarshalIndent(err, "", "  ")
	assert.Nil(marshalErr)
	expected := `{
  "class": "caused by something invalid",
  "xfcc": "mm=hm",
  "metadata": {
    "ecks": "y"
  }
}`
	assert.Equal(expected, string(asBytes))
}

func TestXFCCValidationErrorError(t *testing.T) {
	assert := sdkAssert.New(t)

	c := ex.Class("oh an invalid thing happened")
	var err error = &envoyutil.XFCCValidationError{Class: c}
	assert.Equal(c, err.Error())
}

func TestIsValidationError(t *testing.T) {
	assert := sdkAssert.New(t)

	var err error = ex.New("NOPE")
	assert.False(envoyutil.IsValidationError(err))
	err = &envoyutil.XFCCValidationError{Class: "YEP"}
	assert.True(envoyutil.IsValidationError(err))
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

func TestIsFatalError(t *testing.T) {
	assert := sdkAssert.New(t)

	var err error = ex.New("NOPE")
	assert.False(envoyutil.IsFatalError(err))
	err = &envoyutil.XFCCFatalError{Class: "YEP"}
	assert.True(envoyutil.IsFatalError(err))
}
