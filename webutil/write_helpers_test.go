/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestWriteNoContent(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	res := NewMockResponse(buf)
	assert.Nil(WriteNoContent(res))
	assert.Equal(http.StatusNoContent, res.StatusCode())
	assert.Zero(buf.Len())
}

func TestWriteRawContent(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	res := NewMockResponse(buf)
	assert.Nil(WriteRawContent(res, http.StatusOK, []byte("foo bar baz")))
	assert.Equal(http.StatusOK, res.StatusCode())
	assert.Equal("foo bar baz", buf.String())
}

func TestWriteJSON(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	res := NewMockResponse(buf)
	assert.Nil(WriteJSON(res, http.StatusOK, map[string]interface{}{"foo": "bar"}))
	assert.Equal(http.StatusOK, res.StatusCode())
	assert.Equal("{\"foo\":\"bar\"}\n", buf.String())
}

type flakyWriter struct {
	Err error
}

func (fw flakyWriter) Write(_ []byte) (int, error) {
	return 0, fw.Err
}

func TestWriteJSON_Error(t *testing.T) {
	assert := assert.New(t)

	flaky := flakyWriter{
		Err: fmt.Errorf("flaky error"),
	}
	res := NewMockResponse(flaky)
	err := WriteJSON(res, http.StatusOK, map[string]interface{}{"foo": "bar"})
	assert.NotNil(err)
	assert.Equal("flaky error", err.Error())
	assert.Equal(http.StatusOK, res.StatusCode())
}

func TestWriteJSON_Error_NetOp(t *testing.T) {
	assert := assert.New(t)

	flaky := flakyWriter{
		Err: &net.OpError{
			Op: "test",
		},
	}
	res := NewMockResponse(flaky)
	err := WriteJSON(res, http.StatusOK, map[string]interface{}{"foo": "bar"})
	assert.NotNil(err)
	assert.Equal(ErrNetWrite, ex.ErrClass(err))
	assert.Equal(http.StatusOK, res.StatusCode())
}

type xmltest struct {
	Foo string `xml:"foo"`
}

func TestWriteXML(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	res := NewMockResponse(buf)
	assert.Nil(WriteXML(res, http.StatusOK, xmltest{Foo: "bar"}))
	assert.Equal(http.StatusOK, res.StatusCode())
	assert.Equal("<xmltest><foo>bar</foo></xmltest>", buf.String())
}

func TestWriteXML_Error(t *testing.T) {
	assert := assert.New(t)

	flaky := flakyWriter{
		Err: fmt.Errorf("flaky error"),
	}
	res := NewMockResponse(flaky)
	err := WriteXML(res, http.StatusOK, xmltest{Foo: "bar"})
	assert.NotNil(err)
	assert.Equal("flaky error", err.Error())
	assert.Equal(http.StatusOK, res.StatusCode())
}

func TestWriteXML_Error_NetOp(t *testing.T) {
	assert := assert.New(t)

	flaky := flakyWriter{
		Err: &net.OpError{
			Op: "test",
		},
	}
	res := NewMockResponse(flaky)
	err := WriteXML(res, http.StatusOK, xmltest{Foo: "bar"})
	assert.NotNil(err)
	assert.Equal(ErrNetWrite, ex.ErrClass(err))
	assert.Equal(http.StatusOK, res.StatusCode())
}

func TestDeserializeReaderAsJSON(t *testing.T) {
	assert := assert.New(t)

	contents, err := json.Marshal(map[string]interface{}{"foo": "bar"})
	assert.Nil(err)

	output := make(map[string]interface{})

	assert.Nil(DeserializeReaderAsJSON(&output, io.NopCloser(bytes.NewBuffer(contents))))
	assert.Equal("bar", output["foo"])
}
