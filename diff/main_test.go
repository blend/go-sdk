/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package diff

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/blend/go-sdk/assert"
)

func TestDiffCommonPrefix(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Name	string

		Text1	string
		Text2	string

		Expected	int
	}

	dmp := New()

	for i, tc := range []TestCase{
		{"Null", "abc", "xyz", 0},
		{"Non-null", "1234abcdef", "1234xyz", 4},
		{"Whole", "1234", "1234xyz", 4},
	} {
		actual := dmp.diffCommonPrefix(tc.Text1, tc.Text2)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %s", i, tc.Name))
	}
}

func TestCommonPrefixLength(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Text1	string
		Text2	string

		Expected	int
	}

	for i, tc := range []TestCase{
		{"abc", "xyz", 0},
		{"1234abcdef", "1234xyz", 4},
		{"1234", "1234xyz", 4},
	} {
		actual := commonPrefixLength([]rune(tc.Text1), []rune(tc.Text2))
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %#v", i, tc))
	}
}

func TestDiffCommonSuffix(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Name	string

		Text1	string
		Text2	string

		Expected	int
	}

	dmp := New()

	for i, tc := range []TestCase{
		{"Null", "abc", "xyz", 0},
		{"Non-null", "abcdef1234", "xyz1234", 4},
		{"Whole", "1234", "xyz1234", 4},
	} {
		actual := dmp.diffCommonSuffix(tc.Text1, tc.Text2)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %s", i, tc.Name))
	}
}

func TestCommonSuffixLength(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Text1	string
		Text2	string

		Expected	int
	}

	for i, tc := range []TestCase{
		{"abc", "xyz", 0},
		{"abcdef1234", "xyz1234", 4},
		{"1234", "xyz1234", 4},
		{"123", "a3", 1},
	} {
		actual := commonSuffixLength([]rune(tc.Text1), []rune(tc.Text2))
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %#v", i, tc))
	}
}

func TestDiffCommonOverlap(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Name	string

		Text1	string
		Text2	string

		Expected	int
	}

	dmp := New()

	for i, tc := range []TestCase{
		{"Null", "", "abcd", 0},
		{"Whole", "abc", "abcd", 3},
		{"Null", "123456", "abcd", 0},
		{"Null", "123456xxx", "xxxabcd", 3},
		// Some overly clever languages (C#) may treat ligatures as equal to their component letters, e.g. U+FB01 == 'fi'
		{"Unicode", "fi", "\ufb01i", 0},
	} {
		actual := dmp.diffCommonOverlap(tc.Text1, tc.Text2)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %s", i, tc.Name))
	}
}

func TestDiffHalfMatch(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Text1	string
		Text2	string

		Expected	[]string
	}

	dmp := New()
	dmp.Timeout = 1

	for i, tc := range []TestCase{
		// No match
		{"1234567890", "abcdef", nil},
		{"12345", "23", nil},

		// Single Match
		{"1234567890", "a345678z", []string{"12", "90", "a", "z", "345678"}},
		{"a345678z", "1234567890", []string{"a", "z", "12", "90", "345678"}},
		{"abc56789z", "1234567890", []string{"abc", "z", "1234", "0", "56789"}},
		{"a23456xyz", "1234567890", []string{"a", "xyz", "1", "7890", "23456"}},

		// Multiple Matches
		{"121231234123451234123121", "a1234123451234z", []string{"12123", "123121", "a", "z", "1234123451234"}},
		{"x-=-=-=-=-=-=-=-=-=-=-=-=", "xx-=-=-=-=-=-=-=", []string{"", "-=-=-=-=-=", "x", "", "x-=-=-=-=-=-=-="}},
		{"-=-=-=-=-=-=-=-=-=-=-=-=y", "-=-=-=-=-=-=-=yy", []string{"-=-=-=-=-=", "", "", "y", "-=-=-=-=-=-=-=y"}},

		// Non-optimal halfmatch, ptimal diff would be -q+x=H-i+e=lloHe+Hu=llo-Hew+y not -qHillo+x=HelloHe-w+Hulloy
		{"qHilloHelloHew", "xHelloHeHulloy", []string{"qHillo", "w", "x", "Hulloy", "HelloHe"}},
	} {
		actual := dmp.DiffHalfMatch(tc.Text1, tc.Text2)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %#v", i, tc))
	}

	dmp.Timeout = 0

	for i, tc := range []TestCase{
		// Optimal no halfmatch
		{"qHilloHelloHew", "xHelloHeHulloy", nil},
	} {
		actual := dmp.DiffHalfMatch(tc.Text1, tc.Text2)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %#v", i, tc))
	}
}

func TestDiffBisectSplit(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Text1	string
		Text2	string
	}

	dmp := New()

	for _, tc := range []TestCase{
		{"STUV\x05WX\x05YZ\x05[", "WĺĻļ\x05YZ\x05ĽľĿŀZ"},
	} {
		diffs := dmp.diffBisectSplit([]rune(tc.Text1),
			[]rune(tc.Text2), 7, 6, time.Now().Add(time.Hour))

		for _, d := range diffs {
			assert.True(utf8.ValidString(d.Text))
		}

		// TODO define the expected outcome
	}
}

func TestDiffLinesToChars(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Text1	string
		Text2	string

		ExpectedChars1	string
		ExpectedChars2	string
		ExpectedLines	[]string
	}

	dmp := New()

	for i, tc := range []TestCase{
		{"", "alpha\r\nbeta\r\n\r\n\r\n", "", "1,2,3,3", []string{"", "alpha\r\n", "beta\r\n", "\r\n"}},
		{"a", "b", "1", "2", []string{"", "a", "b"}},
		// Omit final newline.
		{"alpha\nbeta\nalpha", "", "1,2,3", "", []string{"", "alpha\n", "beta\n", "alpha"}},
	} {
		actualChars1, actualChars2, actualLines := dmp.diffLinesToChars(tc.Text1, tc.Text2)
		assert.Equal(tc.ExpectedChars1, actualChars1, fmt.Sprintf("Test case #%d, %#v", i, tc))
		assert.Equal(tc.ExpectedChars2, actualChars2, fmt.Sprintf("Test case #%d, %#v", i, tc))
		assert.Equal(tc.ExpectedLines, actualLines, fmt.Sprintf("Test case #%d, %#v", i, tc))
	}

	// More than 256 to reveal any 8-bit limitations.
	n := 300
	lineList := []string{
		"",	// Account for the initial empty element of the lines array.
	}
	var charList []string
	for x := 1; x < n+1; x++ {
		lineList = append(lineList, strconv.Itoa(x)+"\n")
		charList = append(charList, strconv.Itoa(x))
	}
	lines := strings.Join(lineList, "")
	chars := strings.Join(charList[:], ",")
	assert.Equal(n, len(strings.Split(chars, ",")))

	actualChars1, actualChars2, actualLines := dmp.diffLinesToChars(lines, "")
	assert.Equal(chars, actualChars1)
	assert.Equal("", actualChars2)
	assert.Equal(lineList, actualLines)
}

func TestDiffCharsToLines(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Diffs	[]Diff
		Lines	[]string

		Expected	[]Diff
	}

	dmp := New()

	for i, tc := range []TestCase{
		{
			Diffs: []Diff{
				{DiffEqual, "1,2,1"},
				{DiffInsert, "2,1,2"},
			},
			Lines:	[]string{"", "alpha\n", "beta\n"},

			Expected: []Diff{
				{DiffEqual, "alpha\nbeta\nalpha\n"},
				{DiffInsert, "beta\nalpha\nbeta\n"},
			},
		},
	} {
		actual := dmp.diffCharsToLines(tc.Diffs, tc.Lines)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %#v", i, tc))
	}

	// More than 256 to reveal any 8-bit limitations.
	n := 300
	lineList := []string{
		"",	// Account for the initial empty element of the lines array.
	}
	charList := []string{}
	for x := 1; x <= n; x++ {
		lineList = append(lineList, strconv.Itoa(x)+"\n")
		charList = append(charList, strconv.Itoa(x))
	}
	assert.Equal(n, len(charList))
	chars := strings.Join(charList[:], ",")

	actual := dmp.diffCharsToLines([]Diff{{DiffDelete, chars}}, lineList)
	assert.Equal([]Diff{{DiffDelete, strings.Join(lineList, "")}}, actual)
}

func TestDiffCleanupMerge(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Name	string

		Diffs	[]Diff

		Expected	[]Diff
	}

	dmp := New()

	for i, tc := range []TestCase{
		{
			"Null case",
			[]Diff{},
			[]Diff{},
		},
		{
			"No Diff case",
			[]Diff{{DiffEqual, "a"}, {DiffDelete, "b"}, {DiffInsert, "c"}},
			[]Diff{{DiffEqual, "a"}, {DiffDelete, "b"}, {DiffInsert, "c"}},
		},
		{
			"Merge equalities",
			[]Diff{{DiffEqual, "a"}, {DiffEqual, "b"}, {DiffEqual, "c"}},
			[]Diff{{DiffEqual, "abc"}},
		},
		{
			"Merge deletions",
			[]Diff{{DiffDelete, "a"}, {DiffDelete, "b"}, {DiffDelete, "c"}},
			[]Diff{{DiffDelete, "abc"}},
		},
		{
			"Merge insertions",
			[]Diff{{DiffInsert, "a"}, {DiffInsert, "b"}, {DiffInsert, "c"}},
			[]Diff{{DiffInsert, "abc"}},
		},
		{
			"Merge interweave",
			[]Diff{{DiffDelete, "a"}, {DiffInsert, "b"}, {DiffDelete, "c"}, {DiffInsert, "d"}, {DiffEqual, "e"}, {DiffEqual, "f"}},
			[]Diff{{DiffDelete, "ac"}, {DiffInsert, "bd"}, {DiffEqual, "ef"}},
		},
		{
			"Prefix and suffix detection",
			[]Diff{{DiffDelete, "a"}, {DiffInsert, "abc"}, {DiffDelete, "dc"}},
			[]Diff{{DiffEqual, "a"}, {DiffDelete, "d"}, {DiffInsert, "b"}, {DiffEqual, "c"}},
		},
		{
			"Prefix and suffix detection with equalities",
			[]Diff{{DiffEqual, "x"}, {DiffDelete, "a"}, {DiffInsert, "abc"}, {DiffDelete, "dc"}, {DiffEqual, "y"}},
			[]Diff{{DiffEqual, "xa"}, {DiffDelete, "d"}, {DiffInsert, "b"}, {DiffEqual, "cy"}},
		},
		{
			"Same test as above but with unicode (\u0101 will appear in diffs with at least 257 unique lines)",
			[]Diff{{DiffEqual, "x"}, {DiffDelete, "\u0101"}, {DiffInsert, "\u0101bc"}, {DiffDelete, "dc"}, {DiffEqual, "y"}},
			[]Diff{{DiffEqual, "x\u0101"}, {DiffDelete, "d"}, {DiffInsert, "b"}, {DiffEqual, "cy"}},
		},
		{
			"Slide edit left",
			[]Diff{{DiffEqual, "a"}, {DiffInsert, "ba"}, {DiffEqual, "c"}},
			[]Diff{{DiffInsert, "ab"}, {DiffEqual, "ac"}},
		},
		{
			"Slide edit right",
			[]Diff{{DiffEqual, "c"}, {DiffInsert, "ab"}, {DiffEqual, "a"}},
			[]Diff{{DiffEqual, "ca"}, {DiffInsert, "ba"}},
		},
		{
			"Slide edit left recursive",
			[]Diff{{DiffEqual, "a"}, {DiffDelete, "b"}, {DiffEqual, "c"}, {DiffDelete, "ac"}, {DiffEqual, "x"}},
			[]Diff{{DiffDelete, "abc"}, {DiffEqual, "acx"}},
		},
		{
			"Slide edit right recursive",
			[]Diff{{DiffEqual, "x"}, {DiffDelete, "ca"}, {DiffEqual, "c"}, {DiffDelete, "b"}, {DiffEqual, "a"}},
			[]Diff{{DiffEqual, "xca"}, {DiffDelete, "cba"}},
		},
	} {
		actual := dmp.diffCleanupMerge(tc.Diffs)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %s", i, tc.Name))
	}
}

func TestDiffCleanupSemanticLossless(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Name	string

		Diffs	[]Diff

		Expected	[]Diff
	}

	dmp := New()

	for i, tc := range []TestCase{
		{
			"Null case",
			[]Diff{},
			[]Diff{},
		},
		{
			"Blank lines",
			[]Diff{
				{DiffEqual, "AAA\r\n\r\nBBB"},
				{DiffInsert, "\r\nDDD\r\n\r\nBBB"},
				{DiffEqual, "\r\nEEE"},
			},
			[]Diff{
				{DiffEqual, "AAA\r\n\r\n"},
				{DiffInsert, "BBB\r\nDDD\r\n\r\n"},
				{DiffEqual, "BBB\r\nEEE"},
			},
		},
		{
			"Line boundaries",
			[]Diff{
				{DiffEqual, "AAA\r\nBBB"},
				{DiffInsert, " DDD\r\nBBB"},
				{DiffEqual, " EEE"},
			},
			[]Diff{
				{DiffEqual, "AAA\r\n"},
				{DiffInsert, "BBB DDD\r\n"},
				{DiffEqual, "BBB EEE"},
			},
		},
		{
			"Word boundaries",
			[]Diff{
				{DiffEqual, "The c"},
				{DiffInsert, "ow and the c"},
				{DiffEqual, "at."},
			},
			[]Diff{
				{DiffEqual, "The "},
				{DiffInsert, "cow and the "},
				{DiffEqual, "cat."},
			},
		},
		{
			"Alphanumeric boundaries",
			[]Diff{
				{DiffEqual, "The-c"},
				{DiffInsert, "ow-and-the-c"},
				{DiffEqual, "at."},
			},
			[]Diff{
				{DiffEqual, "The-"},
				{DiffInsert, "cow-and-the-"},
				{DiffEqual, "cat."},
			},
		},
		{
			"Hitting the start",
			[]Diff{
				{DiffEqual, "a"},
				{DiffDelete, "a"},
				{DiffEqual, "ax"},
			},
			[]Diff{
				{DiffDelete, "a"},
				{DiffEqual, "aax"},
			},
		},
		{
			"Hitting the end",
			[]Diff{
				{DiffEqual, "xa"},
				{DiffDelete, "a"},
				{DiffEqual, "a"},
			},
			[]Diff{
				{DiffEqual, "xaa"},
				{DiffDelete, "a"},
			},
		},
		{
			"Sentence boundaries",
			[]Diff{
				{DiffEqual, "The xxx. The "},
				{DiffInsert, "zzz. The "},
				{DiffEqual, "yyy."},
			},
			[]Diff{
				{DiffEqual, "The xxx."},
				{DiffInsert, " The zzz."},
				{DiffEqual, " The yyy."},
			},
		},
		{
			"UTF-8 strings",
			[]Diff{
				{DiffEqual, "The ♕. The "},
				{DiffInsert, "♔. The "},
				{DiffEqual, "♖."},
			},
			[]Diff{
				{DiffEqual, "The ♕."},
				{DiffInsert, " The ♔."},
				{DiffEqual, " The ♖."},
			},
		},
		{
			"Rune boundaries",
			[]Diff{
				{DiffEqual, "♕♕"},
				{DiffInsert, "♔♔"},
				{DiffEqual, "♖♖"},
			},
			[]Diff{
				{DiffEqual, "♕♕"},
				{DiffInsert, "♔♔"},
				{DiffEqual, "♖♖"},
			},
		},
	} {
		actual := dmp.diffCleanupSemanticLossless(tc.Diffs)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %s", i, tc.Name))
	}
}

func TestDiffCleanupSemantic(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Name	string

		Diffs	[]Diff

		Expected	[]Diff
	}

	dmp := New()

	for i, tc := range []TestCase{
		{
			"Null case",
			[]Diff{},
			[]Diff{},
		},
		{
			"No elimination #1",
			[]Diff{
				{DiffDelete, "ab"},
				{DiffInsert, "cd"},
				{DiffEqual, "12"},
				{DiffDelete, "e"},
			},
			[]Diff{
				{DiffDelete, "ab"},
				{DiffInsert, "cd"},
				{DiffEqual, "12"},
				{DiffDelete, "e"},
			},
		},
		{
			"No elimination #2",
			[]Diff{
				{DiffDelete, "abc"},
				{DiffInsert, "ABC"},
				{DiffEqual, "1234"},
				{DiffDelete, "wxyz"},
			},
			[]Diff{
				{DiffDelete, "abc"},
				{DiffInsert, "ABC"},
				{DiffEqual, "1234"},
				{DiffDelete, "wxyz"},
			},
		},
		{
			"No elimination #3",
			[]Diff{
				{DiffEqual, "2016-09-01T03:07:1"},
				{DiffInsert, "5.15"},
				{DiffEqual, "4"},
				{DiffDelete, "."},
				{DiffEqual, "80"},
				{DiffInsert, "0"},
				{DiffEqual, "78"},
				{DiffDelete, "3074"},
				{DiffEqual, "1Z"},
			},
			[]Diff{
				{DiffEqual, "2016-09-01T03:07:1"},
				{DiffInsert, "5.15"},
				{DiffEqual, "4"},
				{DiffDelete, "."},
				{DiffEqual, "80"},
				{DiffInsert, "0"},
				{DiffEqual, "78"},
				{DiffDelete, "3074"},
				{DiffEqual, "1Z"},
			},
		},
		{
			"Simple elimination",
			[]Diff{
				{DiffDelete, "a"},
				{DiffEqual, "b"},
				{DiffDelete, "c"},
			},
			[]Diff{
				{DiffDelete, "abc"},
				{DiffInsert, "b"},
			},
		},
		{
			"Backpass elimination",
			[]Diff{
				{DiffDelete, "ab"},
				{DiffEqual, "cd"},
				{DiffDelete, "e"},
				{DiffEqual, "f"},
				{DiffInsert, "g"},
			},
			[]Diff{
				{DiffDelete, "abcdef"},
				{DiffInsert, "cdfg"},
			},
		},
		{
			"Multiple eliminations",
			[]Diff{
				{DiffInsert, "1"},
				{DiffEqual, "A"},
				{DiffDelete, "B"},
				{DiffInsert, "2"},
				{DiffEqual, "_"},
				{DiffInsert, "1"},
				{DiffEqual, "A"},
				{DiffDelete, "B"},
				{DiffInsert, "2"},
			},
			[]Diff{
				{DiffDelete, "AB_AB"},
				{DiffInsert, "1A2_1A2"},
			},
		},
		{
			"Word boundaries",
			[]Diff{
				{DiffEqual, "The c"},
				{DiffDelete, "ow and the c"},
				{DiffEqual, "at."},
			},
			[]Diff{
				{DiffEqual, "The "},
				{DiffDelete, "cow and the "},
				{DiffEqual, "cat."},
			},
		},
		{
			"No overlap elimination",
			[]Diff{
				{DiffDelete, "abcxx"},
				{DiffInsert, "xxdef"},
			},
			[]Diff{
				{DiffDelete, "abcxx"},
				{DiffInsert, "xxdef"},
			},
		},
		{
			"Overlap elimination",
			[]Diff{
				{DiffDelete, "abcxxx"},
				{DiffInsert, "xxxdef"},
			},
			[]Diff{
				{DiffDelete, "abc"},
				{DiffEqual, "xxx"},
				{DiffInsert, "def"},
			},
		},
		{
			"Reverse overlap elimination",
			[]Diff{
				{DiffDelete, "xxxabc"},
				{DiffInsert, "defxxx"},
			},
			[]Diff{
				{DiffInsert, "def"},
				{DiffEqual, "xxx"},
				{DiffDelete, "abc"},
			},
		},
		{
			"Two overlap eliminations",
			[]Diff{
				{DiffDelete, "abcd1212"},
				{DiffInsert, "1212efghi"},
				{DiffEqual, "----"},
				{DiffDelete, "A3"},
				{DiffInsert, "3BC"},
			},
			[]Diff{
				{DiffDelete, "abcd"},
				{DiffEqual, "1212"},
				{DiffInsert, "efghi"},
				{DiffEqual, "----"},
				{DiffDelete, "A"},
				{DiffEqual, "3"},
				{DiffInsert, "BC"},
			},
		},
		{
			"Test case for adapting DiffCleanupSemantic to be equal to the Python version #19",
			[]Diff{
				{DiffEqual, "James McCarthy "},
				{DiffDelete, "close to "},
				{DiffEqual, "sign"},
				{DiffDelete, "ing"},
				{DiffInsert, "s"},
				{DiffEqual, " new "},
				{DiffDelete, "E"},
				{DiffInsert, "fi"},
				{DiffEqual, "ve"},
				{DiffInsert, "-yea"},
				{DiffEqual, "r"},
				{DiffDelete, "ton"},
				{DiffEqual, " deal"},
				{DiffInsert, " at Everton"},
			},
			[]Diff{
				{DiffEqual, "James McCarthy "},
				{DiffDelete, "close to "},
				{DiffEqual, "sign"},
				{DiffDelete, "ing"},
				{DiffInsert, "s"},
				{DiffEqual, " new "},
				{DiffInsert, "five-year deal at "},
				{DiffEqual, "Everton"},
				{DiffDelete, " deal"},
			},
		},
		{
			"Taken from python / CPP library",
			[]Diff{
				{DiffInsert, "星球大戰：新的希望 "},
				{DiffEqual, "star wars: "},
				{DiffDelete, "episodio iv - un"},
				{DiffEqual, "a n"},
				{DiffDelete, "u"},
				{DiffEqual, "e"},
				{DiffDelete, "va"},
				{DiffInsert, "w"},
				{DiffEqual, " "},
				{DiffDelete, "es"},
				{DiffInsert, "ho"},
				{DiffEqual, "pe"},
				{DiffDelete, "ranza"},
			},
			[]Diff{
				{DiffInsert, "星球大戰：新的希望 "},
				{DiffEqual, "star wars: "},
				{DiffDelete, "episodio iv - una nueva esperanza"},
				{DiffInsert, "a new hope"},
			},
		},
		{
			"panic",
			[]Diff{
				{DiffInsert, "킬러 인 "},
				{DiffEqual, "리커버리"},
				{DiffDelete, " 보이즈"},
			},
			[]Diff{
				{DiffInsert, "킬러 인 "},
				{DiffEqual, "리커버리"},
				{DiffDelete, " 보이즈"},
			},
		},
	} {
		actual := dmp.diffCleanupSemantic(tc.Diffs)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %s", i, tc.Name))
	}
}

func TestDiffCleanupEfficiency(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Name	string

		Diffs	[]Diff

		Expected	[]Diff
	}

	dmp := New()
	dmp.EditCost = 4

	for i, tc := range []TestCase{
		{
			"Null case",
			[]Diff{},
			[]Diff{},
		},
		{
			"No elimination",
			[]Diff{
				{DiffDelete, "ab"},
				{DiffInsert, "12"},
				{DiffEqual, "wxyz"},
				{DiffDelete, "cd"},
				{DiffInsert, "34"},
			},
			[]Diff{
				{DiffDelete, "ab"},
				{DiffInsert, "12"},
				{DiffEqual, "wxyz"},
				{DiffDelete, "cd"},
				{DiffInsert, "34"},
			},
		},
		{
			"Four-edit elimination",
			[]Diff{
				{DiffDelete, "ab"},
				{DiffInsert, "12"},
				{DiffEqual, "xyz"},
				{DiffDelete, "cd"},
				{DiffInsert, "34"},
			},
			[]Diff{
				{DiffDelete, "abxyzcd"},
				{DiffInsert, "12xyz34"},
			},
		},
		{
			"Three-edit elimination",
			[]Diff{
				{DiffInsert, "12"},
				{DiffEqual, "x"},
				{DiffDelete, "cd"},
				{DiffInsert, "34"},
			},
			[]Diff{
				{DiffDelete, "xcd"},
				{DiffInsert, "12x34"},
			},
		},
		{
			"Backpass elimination",
			[]Diff{
				{DiffDelete, "ab"},
				{DiffInsert, "12"},
				{DiffEqual, "xy"},
				{DiffInsert, "34"},
				{DiffEqual, "z"},
				{DiffDelete, "cd"},
				{DiffInsert, "56"},
			},
			[]Diff{
				{DiffDelete, "abxyzcd"},
				{DiffInsert, "12xy34z56"},
			},
		},
	} {
		actual := dmp.diffCleanupEfficiency(tc.Diffs)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %s", i, tc.Name))
	}

	dmp.EditCost = 5

	for i, tc := range []TestCase{
		{
			"High cost elimination",
			[]Diff{
				{DiffDelete, "ab"},
				{DiffInsert, "12"},
				{DiffEqual, "wxyz"},
				{DiffDelete, "cd"},
				{DiffInsert, "34"},
			},
			[]Diff{
				{DiffDelete, "abwxyzcd"},
				{DiffInsert, "12wxyz34"},
			},
		},
	} {
		actual := dmp.diffCleanupEfficiency(tc.Diffs)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %s", i, tc.Name))
	}
}

func Test_PrettyHTML(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Diffs	[]Diff

		Expected	string
	}

	for i, tc := range []TestCase{
		{
			Diffs: []Diff{
				{DiffEqual, "a\n"},
				{DiffDelete, "<B>b</B>"},
				{DiffInsert, "c&d"},
			},

			Expected:	"<span>a&para;<br></span><del style=\"background:#ffe6e6;\">&lt;B&gt;b&lt;/B&gt;</del><ins style=\"background:#e6ffe6;\">c&amp;d</ins>",
		},
	} {
		actual := PrettyHTML(tc.Diffs)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %#v", i, tc))
	}
}

func Test_PrettyText(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Diffs	[]Diff

		Expected	string
	}

	for i, tc := range []TestCase{
		{
			Diffs: []Diff{
				{DiffEqual, "a\n"},
				{DiffDelete, "<B>b</B>"},
				{DiffInsert, "c&d"},
			},

			Expected:	"a\n\x1b[31m<B>b</B>\x1b[0m\x1b[32mc&d\x1b[0m",
		},
	} {
		actual := PrettyText(tc.Diffs)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %#v", i, tc))
	}
}

func Test_Text1_2(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Diffs	[]Diff

		ExpectedText1	string
		ExpectedText2	string
	}

	for i, tc := range []TestCase{
		{
			Diffs: []Diff{
				{DiffEqual, "jump"},
				{DiffDelete, "s"},
				{DiffInsert, "ed"},
				{DiffEqual, " over "},
				{DiffDelete, "the"},
				{DiffInsert, "a"},
				{DiffEqual, " lazy"},
			},

			ExpectedText1:	"jumps over the lazy",
			ExpectedText2:	"jumped over a lazy",
		},
	} {
		actualText1 := Text1(tc.Diffs)
		assert.Equal(tc.ExpectedText1, actualText1, fmt.Sprintf("Test case #%d, %#v", i, tc))

		actualText2 := Text2(tc.Diffs)
		assert.Equal(tc.ExpectedText2, actualText2, fmt.Sprintf("Test case #%d, %#v", i, tc))
	}
}

func TestDiffDelta(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Name	string

		Text	string
		Delta	string

		ErrorMessagePrefix	string
	}

	for i, tc := range []TestCase{
		{"Delta shorter than text", "jumps over the lazyx", "=4\t-1\t+ed\t=6\t-3\t+a\t=5\t+old dog", "Delta length (19) is different from source text length (20)"},
		{"Delta longer than text", "umps over the lazy", "=4\t-1\t+ed\t=6\t-3\t+a\t=5\t+old dog", "Delta length (19) is different from source text length (18)"},
		{"Invalid URL escaping", "", "+%c3%xy", "invalid URL escape \"%xy\""},
		{"Invalid UTF-8 sequence", "", "+%c3xy", "invalid UTF-8 token: \"\\xc3xy\""},
		{"Invalid diff operation", "", "a", "Invalid diff operation in DiffFromDelta: a"},
		{"Invalid diff syntax", "", "-", "strconv.ParseInt: parsing \"\": invalid syntax"},
		{"Negative number in delta", "", "--1", "Negative number in DiffFromDelta: -1"},
		{"Empty case", "", "", ""},
	} {
		diffs, err := FromDelta(tc.Text, tc.Delta)
		msg := fmt.Sprintf("Test case #%d, %s", i, tc.Name)
		if tc.ErrorMessagePrefix == "" {
			assert.Nil(err, msg)
			assert.Nil(diffs, msg)
		} else {
			e := err.Error()
			if strings.HasPrefix(e, tc.ErrorMessagePrefix) {
				e = tc.ErrorMessagePrefix
			}
			assert.Nil(diffs, msg)
			assert.Equal(tc.ErrorMessagePrefix, e, msg)
		}
	}

	// Convert a diff into delta string.
	diffs := []Diff{
		{DiffEqual, "jump"},
		{DiffDelete, "s"},
		{DiffInsert, "ed"},
		{DiffEqual, " over "},
		{DiffDelete, "the"},
		{DiffInsert, "a"},
		{DiffEqual, " lazy"},
		{DiffInsert, "old dog"},
	}
	text1 := Text1(diffs)
	assert.Equal("jumps over the lazy", text1)

	delta := ToDelta(diffs)
	assert.Equal("=4\t-1\t+ed\t=6\t-3\t+a\t=5\t+old dog", delta)

	// Convert delta string into a diff.
	deltaDiffs, err := FromDelta(text1, delta)
	assert.Nil(err)
	assert.Equal(diffs, deltaDiffs)

	// Test deltas with special characters.
	diffs = []Diff{
		{DiffEqual, "\u0680 \x00 \t %"},
		{DiffDelete, "\u0681 \x01 \n ^"},
		{DiffInsert, "\u0682 \x02 \\ |"},
	}
	text1 = Text1(diffs)
	assert.Equal("\u0680 \x00 \t %\u0681 \x01 \n ^", text1)

	// Lowercase, due to UrlEncode uses lower.
	delta = ToDelta(diffs)
	assert.Equal("=7\t-7\t+%DA%82 %02 %5C %7C", delta)

	deltaDiffs, err = FromDelta(text1, delta)
	assert.Equal(diffs, deltaDiffs)
	assert.Nil(err)

	// Verify pool of unchanged characters.
	diffs = []Diff{
		{DiffInsert, "A-Z a-z 0-9 - _ . ! ~ * ' ( ) ; / ? : @ & = + $ , # "},
	}

	delta = ToDelta(diffs)
	assert.Equal("+A-Z a-z 0-9 - _ . ! ~ * ' ( ) ; / ? : @ & = + $ , # ", delta, "Unchanged characters.")

	// Convert delta string into a diff.
	deltaDiffs, err = FromDelta("", delta)
	assert.Equal(diffs, deltaDiffs)
	assert.Nil(err)
}

func TestDiffXIndex(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Name	string

		Diffs		[]Diff
		Location	int

		Expected	int
	}

	dmp := New()

	for i, tc := range []TestCase{
		{"Translation on equality", []Diff{{DiffDelete, "a"}, {DiffInsert, "1234"}, {DiffEqual, "xyz"}}, 2, 5},
		{"Translation on deletion", []Diff{{DiffEqual, "a"}, {DiffDelete, "1234"}, {DiffEqual, "xyz"}}, 3, 1},
	} {
		actual := dmp.diffXIndex(tc.Diffs, tc.Location)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %s", i, tc.Name))
	}
}

func TestDiffLevenshtein(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Name	string

		Diffs	[]Diff

		Expected	int
	}

	for i, tc := range []TestCase{
		{"Levenshtein with trailing equality", []Diff{{DiffDelete, "абв"}, {DiffInsert, "1234"}, {DiffEqual, "эюя"}}, 4},
		{"Levenshtein with leading equality", []Diff{{DiffEqual, "эюя"}, {DiffDelete, "абв"}, {DiffInsert, "1234"}}, 4},
		{"Levenshtein with middle equality", []Diff{{DiffDelete, "абв"}, {DiffEqual, "эюя"}, {DiffInsert, "1234"}}, 7},
	} {
		actual := Levenshtein(tc.Diffs)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %s", i, tc.Name))
	}
}

func TestDiffBisect(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Name	string

		Time	time.Time

		Expected	[]Diff
	}

	dmp := New()

	for i, tc := range []TestCase{
		{
			Name:	"normal",
			Time:	time.Date(9999, time.December, 31, 23, 59, 59, 59, time.UTC),

			Expected: []Diff{
				{DiffDelete, "c"},
				{DiffInsert, "m"},
				{DiffEqual, "a"},
				{DiffDelete, "t"},
				{DiffInsert, "p"},
			},
		},
		{
			Name:	"Negative deadlines count as having infinite time",
			Time:	time.Date(0001, time.January, 01, 00, 00, 00, 00, time.UTC),

			Expected: []Diff{
				{DiffDelete, "c"},
				{DiffInsert, "m"},
				{DiffEqual, "a"},
				{DiffDelete, "t"},
				{DiffInsert, "p"},
			},
		},
		{
			Name:	"Timeout",
			Time:	time.Now().Add(time.Nanosecond),

			Expected: []Diff{
				{DiffDelete, "cat"},
				{DiffInsert, "map"},
			},
		},
	} {
		actual := dmp.diffBisect("cat", "map", tc.Time)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %s", i, tc.Name))
	}

	// Test for invalid UTF-8 sequences
	assert.Equal([]Diff{
		{DiffEqual, "��"},
	}, dmp.diffBisect("\xe0\xe5", "\xe0\xe5", time.Now().Add(time.Minute)))
}

func TestDiff(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Text1	string
		Text2	string

		Expected	[]Diff
	}

	dmp := New()

	// Perform a trivial diff.
	for i, tc := range []TestCase{
		{
			"",
			"",
			nil,
		},
		{
			"abc",
			"abc",
			[]Diff{{DiffEqual, "abc"}},
		},
		{
			"abc",
			"ab123c",
			[]Diff{{DiffEqual, "ab"}, {DiffInsert, "123"}, {DiffEqual, "c"}},
		},
		{
			"a123bc",
			"abc",
			[]Diff{{DiffEqual, "a"}, {DiffDelete, "123"}, {DiffEqual, "bc"}},
		},
		{
			"abc",
			"a123b456c",
			[]Diff{{DiffEqual, "a"}, {DiffInsert, "123"}, {DiffEqual, "b"}, {DiffInsert, "456"}, {DiffEqual, "c"}},
		},
		{
			"a123b456c",
			"abc",
			[]Diff{{DiffEqual, "a"}, {DiffDelete, "123"}, {DiffEqual, "b"}, {DiffDelete, "456"}, {DiffEqual, "c"}},
		},
	} {
		actual := dmp.Diff(tc.Text1, tc.Text2, false)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %#v", i, tc))
	}

	// Perform a real diff and switch off the timeout.
	dmp.Timeout = 0

	for i, tc := range []TestCase{
		{
			"a",
			"b",
			[]Diff{{DiffDelete, "a"}, {DiffInsert, "b"}},
		},
		{
			"Apples are a fruit.",
			"Bananas are also fruit.",
			[]Diff{
				{DiffDelete, "Apple"},
				{DiffInsert, "Banana"},
				{DiffEqual, "s are a"},
				{DiffInsert, "lso"},
				{DiffEqual, " fruit."},
			},
		},
		{
			"ax\t",
			"\u0680x\u0000",
			[]Diff{
				{DiffDelete, "a"},
				{DiffInsert, "\u0680"},
				{DiffEqual, "x"},
				{DiffDelete, "\t"},
				{DiffInsert, "\u0000"},
			},
		},
		{
			"1ayb2",
			"abxab",
			[]Diff{
				{DiffDelete, "1"},
				{DiffEqual, "a"},
				{DiffDelete, "y"},
				{DiffEqual, "b"},
				{DiffDelete, "2"},
				{DiffInsert, "xab"},
			},
		},
		{
			"abcy",
			"xaxcxabc",
			[]Diff{
				{DiffInsert, "xaxcx"},
				{DiffEqual, "abc"}, {DiffDelete, "y"},
			},
		},
		{
			"ABCDa=bcd=efghijklmnopqrsEFGHIJKLMNOefg",
			"a-bcd-efghijklmnopqrs",
			[]Diff{
				{DiffDelete, "ABCD"},
				{DiffEqual, "a"},
				{DiffDelete, "="},
				{DiffInsert, "-"},
				{DiffEqual, "bcd"},
				{DiffDelete, "="},
				{DiffInsert, "-"},
				{DiffEqual, "efghijklmnopqrs"},
				{DiffDelete, "EFGHIJKLMNOefg"},
			},
		},
		{
			"a [[Pennsylvania]] and [[New",
			" and [[Pennsylvania]]",
			[]Diff{
				{DiffInsert, " "},
				{DiffEqual, "a"},
				{DiffInsert, "nd"},
				{DiffEqual, " [[Pennsylvania]]"},
				{DiffDelete, " and [[New"},
			},
		},
	} {
		actual := dmp.Diff(tc.Text1, tc.Text2, false)
		assert.Equal(tc.Expected, actual, fmt.Sprintf("Test case #%d, %#v", i, tc))
	}

	// Test for invalid UTF-8 sequences
	assert.Equal([]Diff{
		{DiffDelete, "��"},
	}, dmp.Diff("\xe0\xe5", "", false))
}

func TestDiffMainWithTimeout(t *testing.T) {
	assert := assert.New(t)

	dmp := New()
	dmp.Timeout = 200 * time.Millisecond

	a := "`Twas brillig, and the slithy toves\nDid gyre and gimble in the wabe:\nAll mimsy were the borogoves,\nAnd the mome raths outgrabe.\n"
	b := "I am the very model of a modern major general,\nI've information vegetable, animal, and mineral,\nI know the kings of England, and I quote the fights historical,\nFrom Marathon to Waterloo, in order categorical.\n"
	// Increase the text lengths by 1024 times to ensure a timeout.
	for x := 0; x < 13; x++ {
		a = a + a
		b = b + b
	}

	startTime := time.Now()
	dmp.Diff(a, b, true)
	endTime := time.Now()

	delta := endTime.Sub(startTime)

	// Test that we took at least the timeout period.
	assert.True(delta >= dmp.Timeout, fmt.Sprintf("%v !>= %v", delta, dmp.Timeout))

	// Test that we didn't take forever (be very forgiving). Theoretically this test could fail very occasionally if the OS task swaps or locks up for a second at the wrong moment.
	assert.True(delta < (dmp.Timeout*100), fmt.Sprintf("%v !< %v", delta, dmp.Timeout*100))
}

func TestDiffMainWithCheckLines(t *testing.T) {
	assert := assert.New(t)

	type TestCase struct {
		Text1	string
		Text2	string
	}

	dmp := New()
	dmp.Timeout = 0

	// Test cases must be at least 100 chars long to pass the cutoff.
	for i, tc := range []TestCase{
		{
			"1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n",
			"abcdefghij\nabcdefghij\nabcdefghij\nabcdefghij\nabcdefghij\nabcdefghij\nabcdefghij\nabcdefghij\nabcdefghij\nabcdefghij\nabcdefghij\nabcdefghij\nabcdefghij\n",
		},
		{
			"1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			"abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghij",
		},
		{
			"1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n1234567890\n",
			"abcdefghij\n1234567890\n1234567890\n1234567890\nabcdefghij\n1234567890\n1234567890\n1234567890\nabcdefghij\n1234567890\n1234567890\n1234567890\nabcdefghij\n",
		},
	} {
		resultWithoutCheckLines := dmp.Diff(tc.Text1, tc.Text2, false)
		resultWithCheckLines := dmp.Diff(tc.Text1, tc.Text2, true)

		// TODO this fails for the third test case, why?
		if i != 2 {
			assert.Equal(resultWithoutCheckLines, resultWithCheckLines, fmt.Sprintf("Test case #%d, %#v", i, tc))
		}
		assert.Equal(diffRebuildTexts(resultWithoutCheckLines), diffRebuildTexts(resultWithCheckLines), fmt.Sprintf("Test case #%d, %#v", i, tc))
	}
}

func pretty(diffs []Diff) string {
	var w bytes.Buffer

	for i, diff := range diffs {
		_, _ = w.WriteString(fmt.Sprintf("%v. ", i))

		switch diff.Type {
		case DiffInsert:
			_, _ = w.WriteString("DiffIns")
		case DiffDelete:
			_, _ = w.WriteString("DiffDel")
		case DiffEqual:
			_, _ = w.WriteString("DiffEql")
		default:
			_, _ = w.WriteString("Unknown")
		}

		_, _ = w.WriteString(fmt.Sprintf(": %v\n", diff.Text))
	}

	return w.String()
}

func diffRebuildTexts(diffs []Diff) []string {
	texts := []string{"", ""}

	for _, d := range diffs {
		if d.Type != DiffInsert {
			texts[0] += d.Text
		}
		if d.Type != DiffDelete {
			texts[1] += d.Text
		}
	}

	return texts
}
