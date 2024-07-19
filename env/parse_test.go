/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package env

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestDelimitedString(t *testing.T) {
	assert := assert.New(t)
	testVars := make(Vars)
	testVars["var_1"] = "val_1"
	res := testVars.DelimitedString(PairDelimiterSemicolon)
	groundTruth := `"var_1"="val_1"`
	assert.Equal(groundTruth, res)

	// Now try with multiple key-val pairs
	testVars["var_2"] = "val_2"
	res = testVars.DelimitedString(PairDelimiterSemicolon)
	groundTruths := []string{`"var_1"="val_1";"var_2"="val_2"`, `"var_2"="val_2";"var_1"="val_1"`}
	t.Log(res)
	assert.True(matchOne(res, groundTruths...))
}

func TestParseGoodInputs(t *testing.T) {
	assert := assert.New(t)

	// Empty string, which is valid
	input := ""
	res, err := ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.Equal(nil, err)

	groundTruth := make(Vars)
	assert.Equal(groundTruth, res)

	// Single valid key-val pair
	input = "var_1=val_1;"
	res, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	assert.Equal(groundTruth, res)

	// Single valid key-val pair with no trailing delimiter
	input = "var_1=val_1"
	res, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	assert.Equal(groundTruth, res)

	// Two valid key-val pairs
	input = "var_1=val_1;var_2=val_2;"
	res, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	groundTruth["var_2"] = "val_2"
	assert.Equal(groundTruth, res)

	// Two valid key-val pairs with arbitrary whitespace
	input = " var_1   = val_1 ;  var_2 =    val_2; "
	res, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.Equal(nil, err)
	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	groundTruth["var_2"] = "val_2"
	assert.Equal(groundTruth, res)

	// Two valid key-val pairs with a quoted string and arbitrary whitespace
	input = "var_1 = val_1; var_2 = \" val_2 \";"
	res, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	groundTruth["var_2"] = " val_2 "
	assert.Equal(groundTruth, res)

	// Two valid key-val pairs with a quoted string and arbitrary whitespace
	// and no trailing separator
	input = `var_1 = val_1; var_2 = " val_2 "`
	res, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	groundTruth["var_2"] = " val_2 "
	assert.Equal(groundTruth, res)

	// Two valid key-val pairs with an escaped quote and arbitrary whitespace
	// and no trailing separator
	input = `  var_1     = val_1  ; var_2 =   \" val_2  `
	res, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	groundTruth["var_2"] = `"val_2`
	assert.Equal(groundTruth, res)

	// Two valid key-val pairs with an escaped quote and arbitrary whitespace
	// and no trailing separator
	input = `var_1 = \=val_1; var_2 = \" val_2 \;  `
	res, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "=val_1"
	groundTruth["var_2"] = `"val_2;`
	assert.Equal(groundTruth, res)

	// two valid key-val pairs where both the key and value for each pair is
	// enclosed in quotes
	input = `"var_1"="val_1";"var_2"="val_2";`
	res, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	groundTruth["var_2"] = "val_2"
	assert.Equal(groundTruth, res)

	// A valid key-val pair consisting of a single quote inside a quoted block
	input = `var_1 = "\""`
	res, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = `"`
	assert.Equal(groundTruth, res)
}

func TestParseBadInputs(t *testing.T) {
	assert := assert.New(t)

	input := "="
	_, err := ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = ";"
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = `\;`
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = "=;"
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = "=some_val;"
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = ";=some_val;"
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = `;\=some_val;`
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = "some"
	res, err := ParsePairDelimiter(input, PairDelimiterSemicolon)
	t.Log(res)
	assert.NotEqual(nil, err)

	input = `some\=val`
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = `key = "`
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = `key "= "`
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = `key \"= "`
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = `key "= \""`
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = `var_1 = =val_1; var_2 = \" val_2 \;  `
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = `var_1 \= val_1; var_2 = " val_2 ";`
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = `var_1 = val_1; var_2 = \" val_2 ";`
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = `var_1 = =val_1; var_1 = \" val_2 \;  `
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = `var_1 \= val_1; var_1 = " val_2 ";`
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)

	input = `var_1 = val_1; var_1 = \" val_2 ";`
	_, err = ParsePairDelimiter(input, PairDelimiterSemicolon)
	assert.NotEqual(nil, err)
}

// TestParseAndBack is an integration test to sanity check that our
// serialization/deserialization methods are at least consistent with each
// other
func TestParseAndBack(t *testing.T) {
	assert := assert.New(t)
	delimiter := PairDelimiterSemicolon

	// Single valid key-val pair
	groundTruth := make(Vars)
	groundTruth["var_1"] = "val_1"
	serialized := groundTruth.DelimitedString(delimiter)
	res, err := ParsePairDelimiter(serialized, PairDelimiterSemicolon)
	assert.Equal(nil, err)
	assert.Equal(groundTruth, res)

	// Single valid key-val pair with no trailing delimiter
	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	groundTruth["var_2"] = "val_2"
	serialized = groundTruth.DelimitedString(delimiter)
	res, _ = ParsePairDelimiter(serialized, PairDelimiterSemicolon)
	assert.Equal(groundTruth, res)

	// Single valid key-val pair with no trailing delimiter
	groundTruth = make(Vars)
	groundTruth[`"`] = "val_1"
	groundTruth[`=`] = "val_2"
	serialized = groundTruth.DelimitedString(delimiter)
	res, _ = ParsePairDelimiter(serialized, PairDelimiterSemicolon)
	assert.Equal(groundTruth, res)

	// More special characters
	groundTruth = make(Vars)
	groundTruth[`\"`] = "val_1"
	groundTruth[`=`] = "val_2"
	serialized = groundTruth.DelimitedString(delimiter)
	res, _ = ParsePairDelimiter(serialized, PairDelimiterSemicolon)
	assert.Equal(groundTruth, res)

	groundTruth = make(Vars)
	groundTruth[`"val_1"="val_2";`] = `"what;a\tricky=value!`
	groundTruth[`=`] = "val_2"
	serialized = groundTruth.DelimitedString(delimiter)
	res, _ = ParsePairDelimiter(serialized, PairDelimiterSemicolon)
	assert.Equal(groundTruth, res)
}

func TestEscapeString(t *testing.T) {
	assert := assert.New(t)

	// no escapes
	input := "some test string"
	expected := "some test string"
	res := escapeString(input, PairDelimiterSemicolon)
	assert.Equal(expected, res)

	input = `some \test string`
	expected = `some \\test string`
	res = escapeString(input, PairDelimiterSemicolon)
	assert.Equal(expected, res)

	input = `some \=test string`
	expected = `some \\\=test string`
	res = escapeString(input, PairDelimiterSemicolon)
	assert.Equal(expected, res)

	input = `test; string`
	expected = `test\; string`
	res = escapeString(input, PairDelimiterSemicolon)
	assert.Equal(expected, res)

	input = `test; " string`
	expected = `test\; \" string`
	res = escapeString(input, PairDelimiterSemicolon)
	assert.Equal(expected, res)
}

// matchOne checks to see if the input string is an exact match of a number of
// candidate ground truths
func matchOne(input string, groundTruths ...string) bool {
	for _, groundTruth := range groundTruths {
		if groundTruth == input {
			return true
		}
	}
	return false
}
