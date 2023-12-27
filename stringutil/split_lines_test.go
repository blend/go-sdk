/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSplitLines(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Input		string
		Expected	[]string
	}{
		{"", nil},
		{"\n", nil},
		{"\n\n", nil},
		{"this", []string{"this"}},
		{"this\nthat", []string{"this", "that"}},
		{"this\nthat\n", []string{"this", "that"}},
		{"this\nthat\nthose", []string{"this", "that", "those"}},
		{"this\nthat\nthose\n", []string{"this", "that", "those"}},
		{"this\nthat\n\nthose\n", []string{"this", "that", "those"}},
		{"this\rthat\nthose\n", []string{"this\rthat", "those"}},
		{"this\rthat\rthose\n", []string{"this\rthat\rthose"}},
		{"this\rthat\rthose\r", []string{"this\rthat\rthose\r"}},
		{"this\r\nthat\rthose\r", []string{"this\r", "that\rthose\r"}},
		{"this\r\nthat\r\nthose\r", []string{"this\r", "that\r", "those\r"}},
		{"this\r\nthat\r\nthose\r\n", []string{"this\r", "that\r", "those\r"}},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, SplitLines(tc.Input,
			OptSplitLinesIncludeNewLine(false),
			OptSplitLinesIncludeEmptyLines(false),
		))
	}
}

func TestSplitLinesIncludeNewLine(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Input		string
		Expected	[]string
	}{
		{"", nil},
		{"\n", nil},
		{"\n\n", nil},
		{"this", []string{"this"}},
		{"this\nthat", []string{"this\n", "that"}},
		{"this\nthat\n", []string{"this\n", "that\n"}},
		{"this\nthat\nthose", []string{"this\n", "that\n", "those"}},
		{"this\nthat\nthose\n", []string{"this\n", "that\n", "those\n"}},
		{"this\nthat\n\nthose\n", []string{"this\n", "that\n", "those\n"}},
		{"this\rthat\nthose\n", []string{"this\rthat\n", "those\n"}},
		{"this\rthat\rthose\n", []string{"this\rthat\rthose\n"}},
		{"this\rthat\rthose\r", []string{"this\rthat\rthose\r"}},
		{"this\r\nthat\rthose\r", []string{"this\r\n", "that\rthose\r"}},
		{"this\r\nthat\r\nthose\r", []string{"this\r\n", "that\r\n", "those\r"}},
		{"this\r\nthat\r\nthose\r\n", []string{"this\r\n", "that\r\n", "those\r\n"}},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, SplitLines(tc.Input,
			OptSplitLinesIncludeNewLine(true),
			OptSplitLinesIncludeEmptyLines(false),
		))
	}
}

func TestSplitLinesIncludeEmptyLines(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Input		string
		Expected	[]string
	}{
		{"", nil},
		{"\n", []string{""}},
		{"\n\n", []string{"", ""}},
		{"this", []string{"this"}},
		{"this\nthat", []string{"this", "that"}},
		{"this\nthat\n", []string{"this", "that"}},
		{"this\nthat\nthose", []string{"this", "that", "those"}},
		{"this\nthat\nthose\n", []string{"this", "that", "those"}},
		{"this\nthat\n\nthose\n", []string{"this", "that", "", "those"}},
		{"this\rthat\nthose\n", []string{"this\rthat", "those"}},
		{"this\rthat\rthose\n", []string{"this\rthat\rthose"}},
		{"this\rthat\rthose\r", []string{"this\rthat\rthose\r"}},
		{"this\r\nthat\rthose\r", []string{"this\r", "that\rthose\r"}},
		{"this\r\nthat\r\nthose\r", []string{"this\r", "that\r", "those\r"}},
		{"this\r\nthat\r\nthose\r\n", []string{"this\r", "that\r", "those\r"}},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, SplitLines(tc.Input,
			OptSplitLinesIncludeNewLine(false),
			OptSplitLinesIncludeEmptyLines(true),
		))
	}
}

func TestSplitLinesIncludeAll(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Input		string
		Expected	[]string
	}{
		{"", nil},
		{"\n", []string{"\n"}},
		{"\n\n", []string{"\n", "\n"}},
		{"this", []string{"this"}},
		{"this\nthat", []string{"this\n", "that"}},
		{"this\nthat\n", []string{"this\n", "that\n"}},
		{"this\nthat\nthose", []string{"this\n", "that\n", "those"}},
		{"this\nthat\nthose\n", []string{"this\n", "that\n", "those\n"}},
		{"this\nthat\n\nthose\n", []string{"this\n", "that\n", "\n", "those\n"}},
		{"this\rthat\nthose\n", []string{"this\rthat\n", "those\n"}},
		{"this\rthat\rthose\n", []string{"this\rthat\rthose\n"}},
		{"this\rthat\rthose\r", []string{"this\rthat\rthose\r"}},
		{"this\r\nthat\rthose\r", []string{"this\r\n", "that\rthose\r"}},
		{"this\r\nthat\r\nthose\r", []string{"this\r\n", "that\r\n", "those\r"}},
		{"this\r\nthat\r\nthose\r\n", []string{"this\r\n", "that\r\n", "those\r\n"}},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, SplitLines(tc.Input,
			OptSplitLinesIncludeNewLine(true),
			OptSplitLinesIncludeEmptyLines(true),
		))
	}
}
