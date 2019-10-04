package main

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

type coverProfileTestCase struct {
	BaseDir  string
	FileName string
	Expected string
}

func TestGopath(t *testing.T) {
	assert := assert.New(t)

	assert.Contains(gopath(), "/go")
}

func TestGlob(t *testing.T) {
	assert := assert.New(t)

	assert.True(glob("", ""))
	assert.True(glob("*", "asdf"))
	assert.True(glob("*/testo/*", "asdf/testo/blah"))
	assert.True(glob("*/*", "asdf/testo"))
	assert.False(glob("*/testo/*", "asdf"))
	assert.False(glob("*/*/*/testo", "asdf/testo"))
	assert.True(glob("*/*/*/testo", "asdf/x/x/testo"))
}

func TestGlobAnyMatch(t *testing.T) {
	assert := assert.New(t)

	assert.False(globAnyMatch("", "testo"))
	assert.True(globAnyMatch("*", "testo"))
	assert.False(globAnyMatch("x", "testo"))
}

func TestJoinCoverPath(t *testing.T) {
	testCases := []coverProfileTestCase{
		{
			"/",
			"foo/bar/baz.go",
			"/foo/bar/baz.go",
		},
		{
			"/users/foo/bar",
			"foo/bar/baz.go",
			"/users/foo/bar/baz.go",
		},
		{
			"/users/bailey/workspace/go/src/github.com/blend/go-sdk/",
			"github.com/blend/go-sdk/assert/assert.go",
			"/users/bailey/workspace/go/src/github.com/blend/go-sdk/assert/assert.go",
		},
		{
			"/go/src/git.blendlabs.com/blend/fees",
			"git.blendlabs.com/blend/fees/pkg/fees/fees.go",
			"/go/src/git.blendlabs.com/blend/fees/pkg/fees/fees.go",
		},
	}

	var actual string
	for _, testCase := range testCases {
		actual = joinCoverPath(testCase.BaseDir, testCase.FileName)
		if actual != testCase.Expected {
			t.Errorf("%s does not match %s", actual, testCase.Expected)
		}
	}
}

func TestPackageFilename(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("testo/asdf", packageFilename("testo", "asdf"))
}

func TestEnforceCoverage(t *testing.T) {
	assert := assert.New(t)

	// errors cases
	assert.NotNil(enforceCoverage("./", "asdf"))
	assert.NotNil(enforceCoverage("/usr/lib", "50"))

	writeCoverage("/tmp", "")
	assert.NotNil(enforceCoverage("/tmp", "90"))

	writeCoverage("/tmp", "90")
	assert.NotNil(enforceCoverage("/tmp", "70"))

	writeCoverage("/tmp", "0")
	assert.Nil(enforceCoverage("/tmp", "0"))

	writeCoverage("/tmp", "70")
	assert.Nil(enforceCoverage("/tmp", "90"))
}

func TestExtractCoverage(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("0", extractCoverage(""))
	assert.Equal("50", extractCoverage("coverage: 50% of statements"))
}

func TestParseCoverage(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(94, parseCoverage("94%"))
	assert.Equal(94, parseCoverage("94"))
}

func TestColorCoverage(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("\x1b[32m90.00\x1b[0m", colorCoverage(90))
	assert.Equal("\x1b[33m75.00\x1b[0m", colorCoverage(75))
	assert.Equal("\x1b[31m30.00\x1b[0m", colorCoverage(30))
	assert.Equal("\x1b[90m0.00\x1b[0m", colorCoverage(0))
}

func TestFormatCoverage(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("75.00", formatCoverage(75))
}
