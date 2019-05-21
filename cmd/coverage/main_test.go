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
