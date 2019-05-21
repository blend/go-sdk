package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

const (
	Foo = `
	package foo

	func bar() int {
		return 1
	}
	`
	FooTest = `
	package foo

	import (
		"testing"
	)

	func TestBar(t *testing.T) {
		bar()
	}
	`
)

type FileInfo struct {
	name string
}

func (fi FileInfo) Name() string {
	return fi.name
}

func (fi FileInfo) Size() int64 {
	return 12
}

func (fi FileInfo) Mode() os.FileMode {
	return 0
}

func (fi FileInfo) IsDir() bool {
	return true
}

func (fi FileInfo) ModTime() time.Time {
	return time.Now()
}

func (fi FileInfo) Sys() interface{} {
	return nil
}

func TestGetPackageCoverageBaseCases(t *testing.T) {
	assert := assert.New(t)

	var packageCoverReport string
	var err error

	_, notExist := os.Stat("fake.xml")
	packageCoverReport, err = getPackageCoverage("./", FileInfo{}, notExist)
	assert.Nil(err)
	assert.Equal("", packageCoverReport)

	blah := errors.New("blah")
	packageCoverReport, err = getPackageCoverage("./", FileInfo{}, blah)
	assert.Equal("", packageCoverReport)
	assert.Equal(blah, err)

	packageCoverReport, err = getPackageCoverage("./", FileInfo{}, nil)
	assert.Nil(err)
	assert.Equal("", packageCoverReport)

	packageCoverReport, err = getPackageCoverage("./testo", FileInfo{name: ".git"}, nil)
	assert.Equal(filepath.SkipDir, err)
	assert.Equal("", packageCoverReport)

	packageCoverReport, err = getPackageCoverage("./testo", FileInfo{name: "_hidden"}, nil)
	assert.Equal(filepath.SkipDir, err)
	assert.Equal("", packageCoverReport)

	packageCoverReport, err = getPackageCoverage("./testo", FileInfo{name: "vendor"}, nil)
	assert.Equal(filepath.SkipDir, err)
	assert.Equal("", packageCoverReport)

	packageCoverReport, err = getPackageCoverage("./testo", FileInfo{name: "/usr/lib"}, nil)
	assert.Nil(err)
	assert.Equal("", packageCoverReport)
}

func createFoo(assert *assert.Assertions) string {
	path := filepath.Join(gopath(), "cov_test")
	err := os.Mkdir(path, os.ModePerm)
	assert.Nil(err)
	ioutil.WriteFile(filepath.Join(path, "foo.go"), []byte(Foo), defaultFileFlags)
	ioutil.WriteFile(filepath.Join(path, "foo_test.go"), []byte(FooTest), defaultFileFlags)
	return path
}

func TestGetPackageCoverage(t *testing.T) {
	assert := assert.New(t)

	var packageCoverReport string
	var err error

	path := createFoo(assert)
	defer os.RemoveAll(path)

	// *exclude = "cov_test"
	// *include = ""

	// packageCoverReport, err = getPackageCoverage(path, FileInfo{name: "cov_test"}, nil)
	// assert.Nil(err)
	// assert.Equal("", packageCoverReport)

	// *exclude = ""
	// *include = "asdf_blah"

	// packageCoverReport, err = getPackageCoverage(path, FileInfo{name: "cov_test"}, nil)
	// assert.Nil(err)
	// assert.Equal("", packageCoverReport)

	// *include = "cov_test"
	// *exclude = "*"

	packageCoverReport, err = getPackageCoverage(path, FileInfo{name: "cov_test"}, nil)
	assert.Nil(err)
	assert.Equal(filepath.Join(path, "profile.cov"), packageCoverReport)
}
