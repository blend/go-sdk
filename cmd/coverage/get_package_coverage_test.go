package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
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

func TestGetPackageCoverageInclude(t *testing.T) {
	assert := assert.New(t)

	*include = "testo/"

	dir, _ := os.Getwd()
	packageCoverReport, err := getPackageCoverage(dir, FileInfo{name: "coverage"}, nil)
	assert.Nil(err)
	assert.Equal("", packageCoverReport)
}

func TestGetPackageCoverageExclude(t *testing.T) {
	assert := assert.New(t)

	*exclude = "cmd/coverage/*"

	dir, _ := os.Getwd()
	packageCoverReport, err := getPackageCoverage(dir, FileInfo{name: "coverage"}, nil)
	assert.Nil(err)
	assert.Equal("", packageCoverReport)
}

func TestGetPackageCoverage(t *testing.T) {
	assert := assert.New(t)

	dir, _ := os.Getwd()
	packageCoverReport, err := getPackageCoverage(filepath.Join(dir, "test"), FileInfo{name: "test"}, nil)
	assert.Nil(err)
	assert.Contains(packageCoverReport, "cmd/coverage/test/profile.cov")
}
