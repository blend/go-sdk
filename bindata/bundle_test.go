package bindata

import (
	"bytes"
	"testing"

	"go/parser"
	"go/token"

	"github.com/blend/go-sdk/assert"
)

func TestBundle(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)
	bundle := new(Bundle)
	bundle.PackageName = "bindata"
	err := bundle.Process(buffer, PathConfig{Path: "./testdata", Recursive: true})
	assert.Nil(err)

	assert.NotEmpty(buffer.Bytes())

	assert.Contains(buffer.String(), "package bindata")
	assert.Contains(buffer.String(), "testdata/js/app.js")
	assert.Contains(buffer.String(), "testdata/css/site.css")

	ast, err := parser.ParseFile(token.NewFileSet(), "bindata.go", buffer.Bytes(), parser.ParseComments|parser.AllErrors)
	assert.Nil(err)
	assert.NotNil(ast)
	assert.Len(ast.Imports, 2)
}
