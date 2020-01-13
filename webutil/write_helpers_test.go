package webutil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
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

func TestDeserializeReaderAsJSON(t *testing.T) {
	assert := assert.New(t)

	contents, err := json.Marshal(map[string]interface{}{"foo": "bar"})
	assert.Nil(err)

	output := make(map[string]interface{})

	assert.Nil(DeserializeReaderAsJSON(&output, ioutil.NopCloser(bytes.NewBuffer(contents))))
	assert.Equal("bar", output["foo"])
}
