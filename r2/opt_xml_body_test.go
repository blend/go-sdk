package r2

import (
	"io/ioutil"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptXMLBody(t *testing.T) {
	assert := assert.New(t)

	r := New("http://localhost:8080", OptXMLBody(xmlTestCase{Status: "OK!"}))
	assert.NotNil(r.Body)

	contents, err := ioutil.ReadAll(r.Body)
	assert.Nil(err)
	assert.Equal("<xmlTestCase><status>OK!</status></xmlTestCase>", string(contents))
}
