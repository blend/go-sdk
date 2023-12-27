/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package diff

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// New creates a new MatchPatch object with default parameters.
func New() *MatchPatch {
	// Defaults.
	return &MatchPatch{
		Timeout:		time.Second,
		EditCost:		4,
		MatchThreshold:		0.5,
		MatchDistance:		1000,
		PatchDeleteThreshold:	0.5,
		PatchMargin:		4,
		MatchMaxBits:		32,
	}
}

// MatchPatch holds the configuration for diff-match-patch operations.
type MatchPatch struct {
	// Number of seconds to map a diff before giving up (0 for infinity).
	Timeout	time.Duration
	// Cost of an empty edit operation in terms of edit characters.
	EditCost	int
	// How far to search for a match (0 = exact location, 1000+ = broad match). A match this many characters away from the expected location will add 1.0 to the score (0.0 is a perfect match).
	MatchDistance	int
	// When deleting a large block of text (over ~64 characters), how close do the contents have to be to match the expected contents. (0.0 = perfection, 1.0 = very loose).  Note that MatchThreshold controls how closely the end points of a delete need to match.
	PatchDeleteThreshold	float64
	// Chunk size for context length.
	PatchMargin	int
	// The number of bits in an int.
	MatchMaxBits	int
	// At what point is no match declared (0.0 = perfection, 1.0 = very loose).
	MatchThreshold	float64
}

// Diff finds the differences between two texts.
//
// If an invalid UTF-8 sequence is encountered, it will be replaced by the Unicode replacement character.
//
// `checklines` indicates if we should do a line level diff, or treat the text as an atomic unit.
func (dmp *MatchPatch) Diff(text1, text2 string, checklines bool) []Diff {
	return dmp.DiffRunes([]rune(text1), []rune(text2), checklines)
}

// DiffRunes finds the differences between two rune sequences.
//
// If an invalid UTF-8 sequence is encountered, it will be replaced by the Unicode replacement character.
//
// `checklines` indicates if we should do a line level diff, or treat the text as an atomic unit.
func (dmp *MatchPatch) DiffRunes(text1, text2 []rune, checklines bool) []Diff {
	var deadline time.Time
	if dmp.Timeout > 0 {
		deadline = time.Now().Add(dmp.Timeout)
	}
	return dmp.diffMainRunes(text1, text2, checklines, deadline)
}

func (dmp *MatchPatch) diffMainRunes(text1, text2 []rune, checklines bool, deadline time.Time) []Diff {
	if runesEqual(text1, text2) {
		var diffs []Diff
		if len(text1) > 0 {
			diffs = append(diffs, Diff{DiffEqual, string(text1)})
		}
		return diffs
	}

	// Trim off common prefix (speedup).
	commonlength := commonPrefixLength(text1, text2)
	commonprefix := text1[:commonlength]
	text1 = text1[commonlength:]
	text2 = text2[commonlength:]

	// Trim off common suffix (speedup).
	commonlength = commonSuffixLength(text1, text2)
	commonsuffix := text1[len(text1)-commonlength:]
	text1 = text1[:len(text1)-commonlength]
	text2 = text2[:len(text2)-commonlength]

	// Compute the diff on the middle block.
	diffs := dmp.diffCompute(text1, text2, checklines, deadline)

	// Restore the prefix and suffix.
	if len(commonprefix) != 0 {
		diffs = append([]Diff{{DiffEqual, string(commonprefix)}}, diffs...)
	}
	if len(commonsuffix) != 0 {
		diffs = append(diffs, Diff{DiffEqual, string(commonsuffix)})
	}

	return dmp.diffCleanupMerge(diffs)
}

// diffCompute finds the differences between two rune slices.  Assumes that the texts do not have any common prefix or suffix.
func (dmp *MatchPatch) diffCompute(text1, text2 []rune, checklines bool, deadline time.Time) []Diff {
	diffs := []Diff{}
	if len(text1) == 0 {
		// Just add some text (speedup).
		return append(diffs, Diff{DiffInsert, string(text2)})
	} else if len(text2) == 0 {
		// Just delete some text (speedup).
		return append(diffs, Diff{DiffDelete, string(text1)})
	}

	var longtext, shorttext []rune
	if len(text1) > len(text2) {
		longtext = text1
		shorttext = text2
	} else {
		longtext = text2
		shorttext = text1
	}

	if i := runesIndex(longtext, shorttext); i != -1 {
		op := DiffInsert
		// Swap insertions for deletions if diff is reversed.
		if len(text1) > len(text2) {
			op = DiffDelete
		}
		// Shorter text is inside the longer text (speedup).
		return []Diff{
			{op, string(longtext[:i])},
			{DiffEqual, string(shorttext)},
			{op, string(longtext[i+len(shorttext):])},
		}
	} else if len(shorttext) == 1 {
		// Single character string.
		// After the previous speedup, the character can't be an equality.
		return []Diff{
			{DiffDelete, string(text1)},
			{DiffInsert, string(text2)},
		}
		// Check to see if the problem can be split in two.
	} else if hm := dmp.diffHalfMatch(text1, text2); hm != nil {
		// A half-match was found, sort out the return data.
		text1A := hm[0]
		text1B := hm[1]
		text2A := hm[2]
		text2B := hm[3]
		midCommon := hm[4]
		// Send both pairs off for separate processing.
		diffsA := dmp.diffMainRunes(text1A, text2A, checklines, deadline)
		diffsB := dmp.diffMainRunes(text1B, text2B, checklines, deadline)
		// Merge the results.
		diffs := diffsA
		diffs = append(diffs, Diff{DiffEqual, string(midCommon)})
		diffs = append(diffs, diffsB...)
		return diffs
	} else if checklines && len(text1) > 100 && len(text2) > 100 {
		return dmp.diffLineMode(text1, text2, deadline)
	}
	return dmp.diffBisectRunes(text1, text2, deadline)
}

// diffLineMode does a quick line-level diff on both []runes, then rediff the parts for greater accuracy. This speedup can produce non-minimal diffs.
func (dmp *MatchPatch) diffLineMode(text1, text2 []rune, deadline time.Time) []Diff {
	// Scan the text on a line-by-line basis first.
	text1, text2, linearray := dmp.diffLinesToRunes(string(text1), string(text2))

	diffs := dmp.diffMainRunes(text1, text2, false, deadline)

	// Convert the diff back to original text.
	diffs = dmp.diffCharsToLines(diffs, linearray)
	// Eliminate freak matches (e.g. blank lines)
	diffs = dmp.diffCleanupSemantic(diffs)

	// Rediff any replacement blocks, this time character-by-character.
	// Add a dummy entry at the end.
	diffs = append(diffs, Diff{DiffEqual, ""})

	pointer := 0
	countDelete := 0
	countInsert := 0

	// NOTE: Rune slices are slower than using strings in this case.
	textDelete := ""
	textInsert := ""

	for pointer < len(diffs) {
		switch diffs[pointer].Type {
		case DiffInsert:
			countInsert++
			textInsert += diffs[pointer].Text
		case DiffDelete:
			countDelete++
			textDelete += diffs[pointer].Text
		case DiffEqual:
			// Upon reaching an equality, check for prior redundancies.
			if countDelete >= 1 && countInsert >= 1 {
				// Delete the offending records and add the merged ones.
				diffs = splice(diffs, pointer-countDelete-countInsert,
					countDelete+countInsert)

				pointer = pointer - countDelete - countInsert
				a := dmp.diffMainRunes([]rune(textDelete), []rune(textInsert), false, deadline)
				for j := len(a) - 1; j >= 0; j-- {
					diffs = splice(diffs, pointer, 0, a[j])
				}
				pointer = pointer + len(a)
			}

			countInsert = 0
			countDelete = 0
			textDelete = ""
			textInsert = ""
		}
		pointer++
	}

	return diffs[:len(diffs)-1]	// Remove the dummy entry at the end.
}

// diffBisect finds the 'middle snake' of a diff, split the problem in two and return the recursively constructed diff.
// If an invalid UTF-8 sequence is encountered, it will be replaced by the Unicode replacement character.
// See Myers 1986 paper: An O(ND) Difference Algorithm and Its Variations.
func (dmp *MatchPatch) diffBisect(text1, text2 string, deadline time.Time) []Diff {
	// Unused in this code, but retained for interface compatibility.
	return dmp.diffBisectRunes([]rune(text1), []rune(text2), deadline)
}

// diffBisect finds the 'middle snake' of a diff, splits the problem in two and returns the recursively constructed diff.
// See Myers's 1986 paper: An O(ND) Difference Algorithm and Its Variations.
func (dmp *MatchPatch) diffBisectRunes(runes1, runes2 []rune, deadline time.Time) []Diff {
	// Cache the text lengths to prevent multiple calls.
	runes1Len, runes2Len := len(runes1), len(runes2)

	maxD := (runes1Len + runes2Len + 1) / 2
	vOffset := maxD
	vLength := 2 * maxD

	v1 := make([]int, vLength)
	v2 := make([]int, vLength)
	for i := range v1 {
		v1[i] = -1
		v2[i] = -1
	}
	v1[vOffset+1] = 0
	v2[vOffset+1] = 0

	delta := runes1Len - runes2Len
	// If the total number of characters is odd, then the front path will collide with the reverse path.
	front := (delta%2 != 0)
	// Offsets for start and end of k loop. Prevents mapping of space beyond the grid.
	k1start := 0
	k1end := 0
	k2start := 0
	k2end := 0
	for d := 0; d < maxD; d++ {
		// Bail out if deadline is reached.
		if !deadline.IsZero() && d%16 == 0 && time.Now().After(deadline) {
			break
		}

		// Walk the front path one step.
		for k1 := -d + k1start; k1 <= d-k1end; k1 += 2 {
			k1Offset := vOffset + k1
			var x1 int

			if k1 == -d || (k1 != d && v1[k1Offset-1] < v1[k1Offset+1]) {
				x1 = v1[k1Offset+1]
			} else {
				x1 = v1[k1Offset-1] + 1
			}

			y1 := x1 - k1
			for x1 < runes1Len && y1 < runes2Len {
				if runes1[x1] != runes2[y1] {
					break
				}
				x1++
				y1++
			}
			v1[k1Offset] = x1
			if x1 > runes1Len {
				// Ran off the right of the graph.
				k1end += 2
			} else if y1 > runes2Len {
				// Ran off the bottom of the graph.
				k1start += 2
			} else if front {
				k2Offset := vOffset + delta - k1
				if k2Offset >= 0 && k2Offset < vLength && v2[k2Offset] != -1 {
					// Mirror x2 onto top-left coordinate system.
					x2 := runes1Len - v2[k2Offset]
					if x1 >= x2 {
						// Overlap detected.
						return dmp.diffBisectSplit(runes1, runes2, x1, y1, deadline)
					}
				}
			}
		}
		// Walk the reverse path one step.
		for k2 := -d + k2start; k2 <= d-k2end; k2 += 2 {
			k2Offset := vOffset + k2
			var x2 int
			if k2 == -d || (k2 != d && v2[k2Offset-1] < v2[k2Offset+1]) {
				x2 = v2[k2Offset+1]
			} else {
				x2 = v2[k2Offset-1] + 1
			}
			var y2 = x2 - k2
			for x2 < runes1Len && y2 < runes2Len {
				if runes1[runes1Len-x2-1] != runes2[runes2Len-y2-1] {
					break
				}
				x2++
				y2++
			}
			v2[k2Offset] = x2
			if x2 > runes1Len {
				// Ran off the left of the graph.
				k2end += 2
			} else if y2 > runes2Len {
				// Ran off the top of the graph.
				k2start += 2
			} else if !front {
				k1Offset := vOffset + delta - k2
				if k1Offset >= 0 && k1Offset < vLength && v1[k1Offset] != -1 {
					x1 := v1[k1Offset]
					y1 := vOffset + x1 - k1Offset
					// Mirror x2 onto top-left coordinate system.
					x2 = runes1Len - x2
					if x1 >= x2 {
						// Overlap detected.
						return dmp.diffBisectSplit(runes1, runes2, x1, y1, deadline)
					}
				}
			}
		}
	}
	// Diff took too long and hit the deadline or number of diffs equals number of characters, no commonality at all.
	return []Diff{
		{DiffDelete, string(runes1)},
		{DiffInsert, string(runes2)},
	}
}

func (dmp *MatchPatch) diffBisectSplit(runes1, runes2 []rune, x, y int,
	deadline time.Time) []Diff {
	runes1a := runes1[:x]
	runes2a := runes2[:y]
	runes1b := runes1[x:]
	runes2b := runes2[y:]

	// Compute both diffs serially.
	diffs := dmp.diffMainRunes(runes1a, runes2a, false, deadline)
	diffsb := dmp.diffMainRunes(runes1b, runes2b, false, deadline)

	return append(diffs, diffsb...)
}

// diffLinesToChars splits two texts into a list of strings, and educes the texts to a string of hashes where each Unicode character represents one line.
// It's slightly faster to call DiffLinesToRunes first, followed by DiffMainRunes.
func (dmp *MatchPatch) diffLinesToChars(text1, text2 string) (string, string, []string) {
	chars1, chars2, lineArray := dmp.diffLinesToStrings(text1, text2)
	return chars1, chars2, lineArray
}

// diffLinesToRunes splits two texts into a list of runes.
func (dmp *MatchPatch) diffLinesToRunes(text1, text2 string) ([]rune, []rune, []string) {
	chars1, chars2, lineArray := dmp.diffLinesToStrings(text1, text2)
	return []rune(chars1), []rune(chars2), lineArray
}

// diffCharsToLines rehydrates the text in a diff from a string of line hashes to real lines of text.
func (dmp *MatchPatch) diffCharsToLines(diffs []Diff, lineArray []string) []Diff {
	hydrated := make([]Diff, 0, len(diffs))
	for _, aDiff := range diffs {
		chars := strings.Split(aDiff.Text, IndexSeparator)
		text := make([]string, len(chars))

		for i, r := range chars {
			i1, err := strconv.Atoi(r)
			if err == nil {
				text[i] = lineArray[i1]
			}
		}

		aDiff.Text = strings.Join(text, "")
		hydrated = append(hydrated, aDiff)
	}
	return hydrated
}

// DiffCommonPrefix determines the common prefix length of two strings.
func (dmp *MatchPatch) diffCommonPrefix(text1, text2 string) int {
	// Unused in this code, but retained for interface compatibility.
	return commonPrefixLength([]rune(text1), []rune(text2))
}

// diffCommonSuffix determines the common suffix length of two strings.
func (dmp *MatchPatch) diffCommonSuffix(text1, text2 string) int {
	// Unused in this code, but retained for interface compatibility.
	return commonSuffixLength([]rune(text1), []rune(text2))
}

// diffCommonOverlap determines if the suffix of one string is the prefix of another.
func (dmp *MatchPatch) diffCommonOverlap(text1 string, text2 string) int {
	// Cache the text lengths to prevent multiple calls.
	text1Length := len(text1)
	text2Length := len(text2)
	// Eliminate the null case.
	if text1Length == 0 || text2Length == 0 {
		return 0
	}
	// Truncate the longer string.
	if text1Length > text2Length {
		text1 = text1[text1Length-text2Length:]
	} else if text1Length < text2Length {
		text2 = text2[0:text1Length]
	}
	textLength := int(math.Min(float64(text1Length), float64(text2Length)))
	// Quick check for the worst case.
	if text1 == text2 {
		return textLength
	}

	// Start by looking for a single character match and increase length until no match is found. Performance analysis: http://neil.fraser.name/news/2010/11/04/
	best := 0
	length := 1
	for {
		pattern := text1[textLength-length:]
		found := strings.Index(text2, pattern)
		if found == -1 {
			break
		}
		length += found
		if found == 0 || text1[textLength-length:] == text2[0:length] {
			best = length
			length++
		}
	}

	return best
}

// DiffHalfMatch checks whether the two texts share a substring which is at least half the length of the longer text. This speedup can produce non-minimal diffs.
func (dmp *MatchPatch) DiffHalfMatch(text1, text2 string) []string {
	// Unused in this code, but retained for interface compatibility.
	runeSlices := dmp.diffHalfMatch([]rune(text1), []rune(text2))
	if runeSlices == nil {
		return nil
	}

	result := make([]string, len(runeSlices))
	for i, r := range runeSlices {
		result[i] = string(r)
	}
	return result
}

func (dmp *MatchPatch) diffHalfMatch(text1, text2 []rune) [][]rune {
	if dmp.Timeout <= 0 {
		// Don't risk returning a non-optimal diff if we have unlimited time.
		return nil
	}

	var longtext, shorttext []rune
	if len(text1) > len(text2) {
		longtext = text1
		shorttext = text2
	} else {
		longtext = text2
		shorttext = text1
	}

	if len(longtext) < 4 || len(shorttext)*2 < len(longtext) {
		return nil	// Pointless.
	}

	// First check if the second quarter is the seed for a half-match.
	hm1 := dmp.diffHalfMatchI(longtext, shorttext, int(float64(len(longtext)+3)/4))

	// Check again based on the third quarter.
	hm2 := dmp.diffHalfMatchI(longtext, shorttext, int(float64(len(longtext)+1)/2))

	if hm1 == nil && hm2 == nil {
		return nil
	}

	var hm [][]rune
	if hm2 == nil {
		hm = hm1
	} else if hm1 == nil {
		hm = hm2
	} else {
		// Both matched. Select the longest.
		if len(hm1[4]) > len(hm2[4]) {
			hm = hm1
		} else {
			hm = hm2
		}
	}

	// A half-match was found, sort out the return data.
	if len(text1) > len(text2) {
		return hm
	}

	return [][]rune{hm[2], hm[3], hm[0], hm[1], hm[4]}
}

// diffHalfMatchI checks if a substring of shorttext exist within longtext such that the substring is at least half the length of longtext?
// Returns a slice containing the prefix of longtext, the suffix of longtext, the prefix of shorttext, the suffix of shorttext and the common middle, or null if there was no match.
func (dmp *MatchPatch) diffHalfMatchI(l, s []rune, i int) [][]rune {
	var bestCommonA []rune
	var bestCommonB []rune
	var bestCommonLen int
	var bestLongtextA []rune
	var bestLongtextB []rune
	var bestShorttextA []rune
	var bestShorttextB []rune

	// Start with a 1/4 length substring at position i as a seed.
	seed := l[i : i+len(l)/4]

	for j := runesIndexOf(s, seed, 0); j != -1; j = runesIndexOf(s, seed, j+1) {
		prefixLength := commonPrefixLength(l[i:], s[j:])
		suffixLength := commonSuffixLength(l[:i], s[:j])

		if bestCommonLen < suffixLength+prefixLength {
			bestCommonA = s[j-suffixLength : j]
			bestCommonB = s[j : j+prefixLength]
			bestCommonLen = len(bestCommonA) + len(bestCommonB)
			bestLongtextA = l[:i-suffixLength]
			bestLongtextB = l[i+prefixLength:]
			bestShorttextA = s[:j-suffixLength]
			bestShorttextB = s[j+prefixLength:]
		}
	}

	if bestCommonLen*2 < len(l) {
		return nil
	}

	return [][]rune{
		bestLongtextA,
		bestLongtextB,
		bestShorttextA,
		bestShorttextB,
		append(bestCommonA, bestCommonB...),
	}
}

// diffCleanupSemantic reduces the number of edits by eliminating semantically trivial equalities.
func (dmp *MatchPatch) diffCleanupSemantic(diffs []Diff) []Diff {
	changes := false
	// Stack of indices where equalities are found.
	equalities := make([]int, 0, len(diffs))

	var lastequality string
	// Always equal to diffs[equalities[equalitiesLength - 1]][1]
	var pointer int	// Index of current position.
	// Number of characters that changed prior to the equality.
	var lengthInsertions1, lengthDeletions1 int
	// Number of characters that changed after the equality.
	var lengthInsertions2, lengthDeletions2 int

	for pointer < len(diffs) {
		if diffs[pointer].Type == DiffEqual {
			// Equality found.
			equalities = append(equalities, pointer)
			lengthInsertions1 = lengthInsertions2
			lengthDeletions1 = lengthDeletions2
			lengthInsertions2 = 0
			lengthDeletions2 = 0
			lastequality = diffs[pointer].Text
		} else {
			// An insertion or deletion.

			if diffs[pointer].Type == DiffInsert {
				lengthInsertions2 += utf8.RuneCountInString(diffs[pointer].Text)
			} else {
				lengthDeletions2 += utf8.RuneCountInString(diffs[pointer].Text)
			}
			// Eliminate an equality that is smaller or equal to the edits on both sides of it.
			difference1 := int(math.Max(float64(lengthInsertions1), float64(lengthDeletions1)))
			difference2 := int(math.Max(float64(lengthInsertions2), float64(lengthDeletions2)))
			if utf8.RuneCountInString(lastequality) > 0 &&
				(utf8.RuneCountInString(lastequality) <= difference1) &&
				(utf8.RuneCountInString(lastequality) <= difference2) {
				// Duplicate record.
				insPoint := equalities[len(equalities)-1]
				diffs = splice(diffs, insPoint, 0, Diff{DiffDelete, lastequality})

				// Change second copy to insert.
				diffs[insPoint+1].Type = DiffInsert
				// Throw away the equality we just deleted.
				equalities = equalities[:len(equalities)-1]

				if len(equalities) > 0 {
					equalities = equalities[:len(equalities)-1]
				}
				pointer = -1
				if len(equalities) > 0 {
					pointer = equalities[len(equalities)-1]
				}

				lengthInsertions1 = 0	// Reset the counters.
				lengthDeletions1 = 0
				lengthInsertions2 = 0
				lengthDeletions2 = 0
				lastequality = ""
				changes = true
			}
		}
		pointer++
	}

	// Normalize the diff.
	if changes {
		diffs = dmp.diffCleanupMerge(diffs)
	}
	diffs = dmp.diffCleanupSemanticLossless(diffs)
	// Find any overlaps between deletions and insertions.
	// e.g: <del>abcxxx</del><ins>xxxdef</ins>
	//   -> <del>abc</del>xxx<ins>def</ins>
	// e.g: <del>xxxabc</del><ins>defxxx</ins>
	//   -> <ins>def</ins>xxx<del>abc</del>
	// Only extract an overlap if it is as big as the edit ahead or behind it.
	pointer = 1
	for pointer < len(diffs) {
		if diffs[pointer-1].Type == DiffDelete &&
			diffs[pointer].Type == DiffInsert {
			deletion := diffs[pointer-1].Text
			insertion := diffs[pointer].Text
			overlapLength1 := dmp.diffCommonOverlap(deletion, insertion)
			overlapLength2 := dmp.diffCommonOverlap(insertion, deletion)
			if overlapLength1 >= overlapLength2 {
				if float64(overlapLength1) >= float64(utf8.RuneCountInString(deletion))/2 ||
					float64(overlapLength1) >= float64(utf8.RuneCountInString(insertion))/2 {

					// Overlap found. Insert an equality and trim the surrounding edits.
					diffs = splice(diffs, pointer, 0, Diff{DiffEqual, insertion[:overlapLength1]})
					diffs[pointer-1].Text =
						deletion[0 : len(deletion)-overlapLength1]
					diffs[pointer+1].Text = insertion[overlapLength1:]
					pointer++
				}
			} else {
				if float64(overlapLength2) >= float64(utf8.RuneCountInString(deletion))/2 ||
					float64(overlapLength2) >= float64(utf8.RuneCountInString(insertion))/2 {
					// Reverse overlap found. Insert an equality and swap and trim the surrounding edits.
					overlap := Diff{DiffEqual, deletion[:overlapLength2]}
					diffs = splice(diffs, pointer, 0, overlap)
					diffs[pointer-1].Type = DiffInsert
					diffs[pointer-1].Text = insertion[0 : len(insertion)-overlapLength2]
					diffs[pointer+1].Type = DiffDelete
					diffs[pointer+1].Text = deletion[overlapLength2:]
					pointer++
				}
			}
			pointer++
		}
		pointer++
	}

	return diffs
}

// Define some regex patterns for matching boundaries.
var (
	nonAlphaNumericRegex	= regexp.MustCompile(`[^a-zA-Z0-9]`)
	whitespaceRegex		= regexp.MustCompile(`\s`)
	linebreakRegex		= regexp.MustCompile(`[\r\n]`)
	blanklineEndRegex	= regexp.MustCompile(`\n\r?\n$`)
	blanklineStartRegex	= regexp.MustCompile(`^\r?\n\r?\n`)
)

// diffCleanupSemanticScore computes a score representing whether the internal boundary falls on logical boundaries.
// Scores range from 6 (best) to 0 (worst). Closure, but does not reference any external variables.
func (dmp *MatchPatch) diffCleanupSemanticScore(one, two string) int {
	if len(one) == 0 || len(two) == 0 {
		// Edges are the best.
		return 6
	}

	// Each port of this function behaves slightly differently due to subtle differences in each language's definition of things like 'whitespace'.  Since this function's purpose is largely cosmetic, the choice has been made to use each language's native features rather than force total conformity.
	rune1, _ := utf8.DecodeLastRuneInString(one)
	rune2, _ := utf8.DecodeRuneInString(two)
	char1 := string(rune1)
	char2 := string(rune2)

	nonAlphaNumeric1 := nonAlphaNumericRegex.MatchString(char1)
	nonAlphaNumeric2 := nonAlphaNumericRegex.MatchString(char2)
	whitespace1 := nonAlphaNumeric1 && whitespaceRegex.MatchString(char1)
	whitespace2 := nonAlphaNumeric2 && whitespaceRegex.MatchString(char2)
	lineBreak1 := whitespace1 && linebreakRegex.MatchString(char1)
	lineBreak2 := whitespace2 && linebreakRegex.MatchString(char2)
	blankLine1 := lineBreak1 && blanklineEndRegex.MatchString(one)
	blankLine2 := lineBreak2 && blanklineEndRegex.MatchString(two)

	if blankLine1 || blankLine2 {
		// Five points for blank lines.
		return 5
	} else if lineBreak1 || lineBreak2 {
		// Four points for line breaks.
		return 4
	} else if nonAlphaNumeric1 && !whitespace1 && whitespace2 {
		// Three points for end of sentences.
		return 3
	} else if whitespace1 || whitespace2 {
		// Two points for whitespace.
		return 2
	} else if nonAlphaNumeric1 || nonAlphaNumeric2 {
		// One point for non-alphanumeric.
		return 1
	}
	return 0
}

// diffCleanupSemanticLossless looks for single edits surrounded on both sides by equalities which can be shifted sideways to align the edit to a word boundary.
// E.g: The c<ins>at c</ins>ame. -> The <ins>cat </ins>came.
func (dmp *MatchPatch) diffCleanupSemanticLossless(diffs []Diff) []Diff {
	pointer := 1

	// Intentionally ignore the first and last element (don't need checking).
	for pointer < len(diffs)-1 {
		if diffs[pointer-1].Type == DiffEqual &&
			diffs[pointer+1].Type == DiffEqual {

			// This is a single edit surrounded by equalities.
			equality1 := diffs[pointer-1].Text
			edit := diffs[pointer].Text
			equality2 := diffs[pointer+1].Text

			// First, shift the edit as far left as possible.
			commonOffset := dmp.diffCommonSuffix(equality1, edit)
			if commonOffset > 0 {
				commonString := edit[len(edit)-commonOffset:]
				equality1 = equality1[0 : len(equality1)-commonOffset]
				edit = commonString + edit[:len(edit)-commonOffset]
				equality2 = commonString + equality2
			}

			// Second, step character by character right, looking for the best fit.
			bestEquality1 := equality1
			bestEdit := edit
			bestEquality2 := equality2
			bestScore := dmp.diffCleanupSemanticScore(equality1, edit) +
				dmp.diffCleanupSemanticScore(edit, equality2)

			for len(edit) != 0 && len(equality2) != 0 {
				_, sz := utf8.DecodeRuneInString(edit)
				if len(equality2) < sz || edit[:sz] != equality2[:sz] {
					break
				}
				equality1 += edit[:sz]
				edit = edit[sz:] + equality2[:sz]
				equality2 = equality2[sz:]
				score := dmp.diffCleanupSemanticScore(equality1, edit) +
					dmp.diffCleanupSemanticScore(edit, equality2)
				// The >= encourages trailing rather than leading whitespace on edits.
				if score >= bestScore {
					bestScore = score
					bestEquality1 = equality1
					bestEdit = edit
					bestEquality2 = equality2
				}
			}

			if diffs[pointer-1].Text != bestEquality1 {
				// We have an improvement, save it back to the diff.
				if len(bestEquality1) != 0 {
					diffs[pointer-1].Text = bestEquality1
				} else {
					diffs = splice(diffs, pointer-1, 1)
					pointer--
				}

				diffs[pointer].Text = bestEdit
				if len(bestEquality2) != 0 {
					diffs[pointer+1].Text = bestEquality2
				} else {
					diffs = append(diffs[:pointer+1], diffs[pointer+2:]...)
					pointer--
				}
			}
		}
		pointer++
	}

	return diffs
}

// diffCleanupEfficiency reduces the number of edits by eliminating operationally trivial equalities.
func (dmp *MatchPatch) diffCleanupEfficiency(diffs []Diff) []Diff {
	changes := false
	// Stack of indices where equalities are found.
	type equality struct {
		data	int
		next	*equality
	}
	var equalities *equality
	// Always equal to equalities[equalitiesLength-1][1]
	lastequality := ""
	pointer := 0	// Index of current position.
	// Is there an insertion operation before the last equality.
	preIns := false
	// Is there a deletion operation before the last equality.
	preDel := false
	// Is there an insertion operation after the last equality.
	postIns := false
	// Is there a deletion operation after the last equality.
	postDel := false
	for pointer < len(diffs) {
		if diffs[pointer].Type == DiffEqual {	// Equality found.
			if len(diffs[pointer].Text) < dmp.EditCost &&
				(postIns || postDel) {
				// Candidate found.
				equalities = &equality{
					data:	pointer,
					next:	equalities,
				}
				preIns = postIns
				preDel = postDel
				lastequality = diffs[pointer].Text
			} else {
				// Not a candidate, and can never become one.
				equalities = nil
				lastequality = ""
			}
			postIns = false
			postDel = false
		} else {	// An insertion or deletion.
			if diffs[pointer].Type == DiffDelete {
				postDel = true
			} else {
				postIns = true
			}

			// Five types to be split:
			// <ins>A</ins><del>B</del>XY<ins>C</ins><del>D</del>
			// <ins>A</ins>X<ins>C</ins><del>D</del>
			// <ins>A</ins><del>B</del>X<ins>C</ins>
			// <ins>A</del>X<ins>C</ins><del>D</del>
			// <ins>A</ins><del>B</del>X<del>C</del>
			var sumPres int
			if preIns {
				sumPres++
			}
			if preDel {
				sumPres++
			}
			if postIns {
				sumPres++
			}
			if postDel {
				sumPres++
			}
			if len(lastequality) > 0 &&
				((preIns && preDel && postIns && postDel) ||
					((len(lastequality) < dmp.EditCost/2) && sumPres == 3)) {

				insPoint := equalities.data

				// Duplicate record.
				diffs = splice(diffs, insPoint, 0, Diff{DiffDelete, lastequality})

				// Change second copy to insert.
				diffs[insPoint+1].Type = DiffInsert
				// Throw away the equality we just deleted.
				equalities = equalities.next
				lastequality = ""

				if preIns && preDel {
					// No changes made which could affect previous entry, keep going.
					postIns = true
					postDel = true
					equalities = nil
				} else {
					if equalities != nil {
						equalities = equalities.next
					}
					if equalities != nil {
						pointer = equalities.data
					} else {
						pointer = -1
					}
					postIns = false
					postDel = false
				}
				changes = true
			}
		}
		pointer++
	}

	if changes {
		diffs = dmp.diffCleanupMerge(diffs)
	}

	return diffs
}

// diffCleanupMerge reorders and merges like edit sections. Merge equalities.
// Any edit section can move as long as it doesn't cross an equality.
func (dmp *MatchPatch) diffCleanupMerge(diffs []Diff) []Diff {
	// Add a dummy entry at the end.
	diffs = append(diffs, Diff{DiffEqual, ""})
	pointer := 0
	countDelete := 0
	countInsert := 0
	commonlength := 0
	textDelete := []rune(nil)
	textInsert := []rune(nil)

	for pointer < len(diffs) {
		switch diffs[pointer].Type {
		case DiffInsert:
			countInsert++
			textInsert = append(textInsert, []rune(diffs[pointer].Text)...)
			pointer++
			break
		case DiffDelete:
			countDelete++
			textDelete = append(textDelete, []rune(diffs[pointer].Text)...)
			pointer++
			break
		case DiffEqual:
			// Upon reaching an equality, check for prior redundancies.
			if countDelete+countInsert > 1 {
				if countDelete != 0 && countInsert != 0 {
					// Factor out any common prefixies.
					commonlength = commonPrefixLength(textInsert, textDelete)
					if commonlength != 0 {
						x := pointer - countDelete - countInsert
						if x > 0 && diffs[x-1].Type == DiffEqual {
							diffs[x-1].Text += string(textInsert[:commonlength])
						} else {
							diffs = append([]Diff{{DiffEqual, string(textInsert[:commonlength])}}, diffs...)
							pointer++
						}
						textInsert = textInsert[commonlength:]
						textDelete = textDelete[commonlength:]
					}
					// Factor out any common suffixies.
					commonlength = commonSuffixLength(textInsert, textDelete)
					if commonlength != 0 {
						insertIndex := len(textInsert) - commonlength
						deleteIndex := len(textDelete) - commonlength
						diffs[pointer].Text = string(textInsert[insertIndex:]) + diffs[pointer].Text
						textInsert = textInsert[:insertIndex]
						textDelete = textDelete[:deleteIndex]
					}
				}
				// Delete the offending records and add the merged ones.
				if countDelete == 0 {
					diffs = splice(diffs, pointer-countInsert,
						countDelete+countInsert,
						Diff{DiffInsert, string(textInsert)})
				} else if countInsert == 0 {
					diffs = splice(diffs, pointer-countDelete,
						countDelete+countInsert,
						Diff{DiffDelete, string(textDelete)})
				} else {
					diffs = splice(diffs, pointer-countDelete-countInsert,
						countDelete+countInsert,
						Diff{DiffDelete, string(textDelete)},
						Diff{DiffInsert, string(textInsert)})
				}

				pointer = pointer - countDelete - countInsert + 1
				if countDelete != 0 {
					pointer++
				}
				if countInsert != 0 {
					pointer++
				}
			} else if pointer != 0 && diffs[pointer-1].Type == DiffEqual {
				// Merge this equality with the previous one.
				diffs[pointer-1].Text += diffs[pointer].Text
				diffs = append(diffs[:pointer], diffs[pointer+1:]...)
			} else {
				pointer++
			}
			countInsert = 0
			countDelete = 0
			textDelete = nil
			textInsert = nil
			break
		}
	}

	if len(diffs[len(diffs)-1].Text) == 0 {
		diffs = diffs[0 : len(diffs)-1]	// Remove the dummy entry at the end.
	}

	// Second pass: look for single edits surrounded on both sides by equalities which can be shifted sideways to eliminate an equality. E.g: A<ins>BA</ins>C -> <ins>AB</ins>AC
	changes := false
	pointer = 1
	// Intentionally ignore the first and last element (don't need checking).
	for pointer < (len(diffs) - 1) {
		if diffs[pointer-1].Type == DiffEqual &&
			diffs[pointer+1].Type == DiffEqual {
			// This is a single edit surrounded by equalities.
			if strings.HasSuffix(diffs[pointer].Text, diffs[pointer-1].Text) {
				// Shift the edit over the previous equality.
				diffs[pointer].Text = diffs[pointer-1].Text +
					diffs[pointer].Text[:len(diffs[pointer].Text)-len(diffs[pointer-1].Text)]
				diffs[pointer+1].Text = diffs[pointer-1].Text + diffs[pointer+1].Text
				diffs = splice(diffs, pointer-1, 1)
				changes = true
			} else if strings.HasPrefix(diffs[pointer].Text, diffs[pointer+1].Text) {
				// Shift the edit over the next equality.
				diffs[pointer-1].Text += diffs[pointer+1].Text
				diffs[pointer].Text =
					diffs[pointer].Text[len(diffs[pointer+1].Text):] + diffs[pointer+1].Text
				diffs = splice(diffs, pointer+1, 1)
				changes = true
			}
		}
		pointer++
	}

	// If shifts were made, the diff needs reordering and another shift sweep.
	if changes {
		diffs = dmp.diffCleanupMerge(diffs)
	}

	return diffs
}

// diffXIndex returns the equivalent location in s2.
func (dmp *MatchPatch) diffXIndex(diffs []Diff, loc int) int {
	chars1 := 0
	chars2 := 0
	lastChars1 := 0
	lastChars2 := 0
	lastDiff := Diff{}
	for i := 0; i < len(diffs); i++ {
		aDiff := diffs[i]
		if aDiff.Type != DiffInsert {
			// Equality or deletion.
			chars1 += len(aDiff.Text)
		}
		if aDiff.Type != DiffDelete {
			// Equality or insertion.
			chars2 += len(aDiff.Text)
		}
		if chars1 > loc {
			// Overshot the location.
			lastDiff = aDiff
			break
		}
		lastChars1 = chars1
		lastChars2 = chars2
	}
	if lastDiff.Type == DiffDelete {
		// The location was deleted.
		return lastChars2
	}
	// Add the remaining character length.
	return lastChars2 + (loc - lastChars1)
}

// diffLinesToStrings splits two texts into a list of strings. Each string represents one line.
func (dmp *MatchPatch) diffLinesToStrings(text1, text2 string) (string, string, []string) {
	// '\x00' is a valid character, but various debuggers don't like it. So we'll insert a junk entry to avoid generating a null character.
	lineArray := []string{""}	// e.g. lineArray[4] == 'Hello\n'

	//Each string has the index of lineArray which it points to
	strIndexArray1 := dmp.diffLinesToStringsMunge(text1, &lineArray)
	strIndexArray2 := dmp.diffLinesToStringsMunge(text2, &lineArray)

	return intArrayToString(strIndexArray1), intArrayToString(strIndexArray2), lineArray
}

// diffLinesToStringsMunge splits a text into an array of strings, and reduces the texts to a []string.
func (dmp *MatchPatch) diffLinesToStringsMunge(text string, lineArray *[]string) []uint32 {
	// Walk the text, pulling out a substring for each line. text.split('\n') would would temporarily double our memory footprint. Modifying text would create many large strings to garbage collect.
	lineHash := map[string]int{}	// e.g. lineHash['Hello\n'] == 4
	lineStart := 0
	lineEnd := -1
	strs := []uint32{}

	for lineEnd < len(text)-1 {
		lineEnd = indexOf(text, "\n", lineStart)

		if lineEnd == -1 {
			lineEnd = len(text) - 1
		}

		line := text[lineStart : lineEnd+1]
		lineStart = lineEnd + 1
		lineValue, ok := lineHash[line]

		if ok {
			strs = append(strs, uint32(lineValue))
		} else {
			*lineArray = append(*lineArray, line)
			lineHash[line] = len(*lineArray) - 1
			strs = append(strs, uint32(len(*lineArray)-1))
		}
	}

	return strs
}

// runesIndex is the equivalent of strings.Index for rune slices.
func runesIndex(r1, r2 []rune) int {
	last := len(r1) - len(r2)
	for i := 0; i <= last; i++ {
		if runesEqual(r1[i:i+len(r2)], r2) {
			return i
		}
	}
	return -1
}

// runesIndexOf returns the index of pattern in target, starting at target[i].
func runesIndexOf(target, pattern []rune, i int) int {
	if i > len(target)-1 {
		return -1
	}
	if i <= 0 {
		return runesIndex(target, pattern)
	}
	ind := runesIndex(target[i:], pattern)
	if ind == -1 {
		return -1
	}
	return ind + i
}

func runesEqual(r1, r2 []rune) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i, c := range r1 {
		if c != r2[i] {
			return false
		}
	}
	return true
}

// indexOf returns the first index of pattern in str, starting at str[i].
func indexOf(str string, pattern string, i int) int {
	if i > len(str)-1 {
		return -1
	}
	if i <= 0 {
		return strings.Index(str, pattern)
	}
	ind := strings.Index(str[i:], pattern)
	if ind == -1 {
		return -1
	}
	return ind + i
}

// lastIndexOf returns the last index of pattern in str, starting at str[i].
func lastIndexOf(str string, pattern string, i int) int {
	if i < 0 {
		return -1
	}
	if i >= len(str) {
		return strings.LastIndex(str, pattern)
	}
	_, size := utf8.DecodeRuneInString(str[i:])
	return strings.LastIndex(str[:i+size], pattern)
}

func intArrayToString(ns []uint32) string {
	if len(ns) == 0 {
		return ""
	}

	indexSeparator := IndexSeparator[0]

	// Appr. 3 chars per num plus the comma.
	b := []byte{}
	for _, n := range ns {
		b = strconv.AppendInt(b, int64(n), 10)
		b = append(b, indexSeparator)
	}
	b = b[:len(b)-1]
	return string(b)
}

// unescaper unescapes selected chars for compatibility with JavaScript's encodeURI.
// In speed critical applications this could be dropped since the receiving application will certainly decode these fine. Note that this function is case-sensitive.  Thus "%3F" would not be unescaped.  But this is ok because it is only called with the output of HttpUtility.UrlEncode which returns lowercase hex. Example: "%3f" -> "?", "%24" -> "$", etc.
var unescaper = strings.NewReplacer(
	"%21", "!", "%7E", "~", "%27", "'",
	"%28", "(", "%29", ")", "%3B", ";",
	"%2F", "/", "%3F", "?", "%3A", ":",
	"%40", "@", "%26", "&", "%3D", "=",
	"%2B", "+", "%24", "$", "%2C", ",", "%23", "#", "%2A", "*")

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// splice removes amount elements from slice at index index, replacing them with elements.
func splice(slice []Diff, index int, amount int, elements ...Diff) []Diff {
	if len(elements) == amount {
		// Easy case: overwrite the relevant items.
		copy(slice[index:], elements)
		return slice
	}
	if len(elements) < amount {
		// Fewer new items than old.
		// Copy in the new items.
		copy(slice[index:], elements)
		// Shift the remaining items left.
		copy(slice[index+len(elements):], slice[index+amount:])
		// Calculate the new end of the slice.
		end := len(slice) - amount + len(elements)
		// Zero stranded elements at end so that they can be garbage collected.
		tail := slice[end:]
		for i := range tail {
			tail[i] = Diff{}
		}
		return slice[:end]
	}
	// More new items than old.
	// Make room in slice for new elements.
	// There's probably an even more efficient way to do this,
	// but this is simple and clear.
	need := len(slice) - amount + len(elements)
	for len(slice) < need {
		slice = append(slice, Diff{})
	}
	// Shift slice elements right to make room for new elements.
	copy(slice[index+len(elements):], slice[index+amount:])
	// Copy in new elements.
	copy(slice[index:], elements)
	return slice
}

// commonPrefixLength returns the length of the common prefix of two rune slices.
func commonPrefixLength(text1, text2 []rune) int {
	// Linear search. See comment in commonSuffixLength.
	n := 0
	for ; n < len(text1) && n < len(text2); n++ {
		if text1[n] != text2[n] {
			return n
		}
	}
	return n
}

// commonSuffixLength returns the length of the common suffix of two rune slices.
func commonSuffixLength(text1, text2 []rune) int {
	// Use linear search rather than the binary search discussed at https://neil.fraser.name/news/2007/10/09/.
	// See discussion at https://github.com/sergi/go-diff/issues/54.
	i1 := len(text1)
	i2 := len(text2)
	for n := 0; ; n++ {
		i1--
		i2--
		if i1 < 0 || i2 < 0 || text1[i1] != text2[i2] {
			return n
		}
	}
}
