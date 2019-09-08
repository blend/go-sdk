package bindata

import (
	"regexp"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFindFiles(t *testing.T) {
	assert := assert.New(t)

	var files []*File
	err := FindFiles("./testdata", nil, func(f *File) {
		files = append(files, f)
	})
	assert.Nil(err)
	assert.Len(files, 2)
	assert.Equal("testdata/css/site.css", files[0].Name)
	assert.Equal("testdata/js/app.js", files[1].Name)
}

func TestFindFilesIgnores(t *testing.T) {
	assert := assert.New(t)

	ignoreCSS := regexp.MustCompile(".css$")

	var files []*File
	err := FindFiles("./testdata", []*regexp.Regexp{ignoreCSS}, func(f *File) {
		files = append(files, f)
	})
	assert.Nil(err)
	assert.Len(files, 1)
	assert.Equal("testdata/js/app.js", files[0].Name)
}
