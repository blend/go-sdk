package bindata

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestBundle(t *testing.T) {
	assert := assert.New(t)

	bundle := new(Bundle)
	bundle.PackageName = "bindata"

	err := bundle.Process(PathConfig{Path: "./testdata", Recursive: true})
	assert.Nil(err)

	assert.Len(bundle.Assets, 2)

	file, ok := bundle.Assets["testdata/css/site.css"]
	assert.True(ok, fmt.Sprintf("%#v", bundle.Assets.Paths()))
	assert.NotNil(file)
	assert.Equal("testdata/css/site.css", file.Name)
	assert.NotZero(file.Modtime)
	assert.NotEmpty(file.MD5)
	assert.NotEmpty(file.Contents)

	file, ok = bundle.Assets["testdata/js/app.js"]
	assert.True(ok)
	assert.NotNil(file)
	assert.Equal("testdata/js/app.js", file.Name)
	assert.NotZero(file.Modtime)
	assert.NotEmpty(file.MD5)
	assert.NotEmpty(file.Contents)

	buffer := new(bytes.Buffer)
	err = bundle.Write(buffer)
	assert.Nil(err)
	assert.NotEmpty(buffer.Bytes())
}
