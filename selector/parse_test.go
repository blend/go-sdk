/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package selector

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestMustParse(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("x == a", MustParse("x==a").String())

	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = ex.New(r)
			}
		}()
		MustParse("x!!")
	}()
	assert.NotNil(err)
}

func TestParseInvalid(t *testing.T) {
	assert := assert.New(t)

	testBadStrings := []string{
		"x=a||y=b",
		"x==a==b",
		"!x=a",
		"x<a",
		"x>1",
		"x>1,z<5",
		"x=",
		"x= ",
		"x=,z= ",
		"x= ,z= ",
		"foo == bar foo",
	}
	var err error
	for _, str := range testBadStrings {
		_, err = Parse(str)
		assert.NotNil(err, str)
	}
}

func TestParseSemiValid(t *testing.T) {
	assert := assert.New(t)

	testGoodStrings := []string{
		"",
		"x=a,y=b,z=c",
		"x!=a,y=b",
		"x",
		"!x",
		"!x,y",
		"!x,y,z",
		"!x,y,!z",
		"!x,y,z,a",
		"!x,y,z,!a",
	}

	var err error
	for _, str := range testGoodStrings {
		_, err = Parse(str)
		assert.Nil(err, str)
	}
}

func TestParseTrailingComma(t *testing.T) {
	assert := assert.New(t)

	testBadStrings := []string{
		",",
		"x,",
		"!x,",
		"foo==bar,",
		"foo!=bar,",
		"foo in (bar,baz),",
		"foo not in (bar,baz),",
	}

	var err error
	for _, str := range testBadStrings {
		_, err = Parse(str)
		assert.NotNil(err, str)
	}
}

func TestParseEquals(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo":	"bar",
		"moo":	"lar",
	}
	invalid := Labels{
		"zoo":	"mar",
		"moo":	"lar",
	}

	selector, err := Parse("foo == bar")
	assert.Nil(err)
	assert.True(selector.Matches(valid))
	assert.False(selector.Matches(invalid))
}

func TestParseNotEquals(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo":	"far",
		"moo":	"lar",
	}
	invalidPresent := Labels{
		"foo":	"bar",
		"moo":	"lar",
	}
	invalidMissing := Labels{
		"zoo":	"mar",
		"moo":	"lar",
	}

	selector, err := Parse("foo != bar")
	assert.Nil(err)
	assert.True(selector.Matches(valid))
	assert.True(selector.Matches(invalidMissing))
	assert.False(selector.Matches(invalidPresent))
}

func TestParseIn(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo":	"far",
		"moo":	"lar",
	}
	valid2 := Labels{
		"foo":	"bar",
		"moo":	"lar",
	}
	invalid := Labels{
		"foo":	"mar",
		"moo":	"lar",
	}
	invalidMissing := Labels{
		"zoo":	"mar",
		"moo":	"lar",
	}

	selector, err := Parse("foo in (bar,far)")
	assert.Nil(err)
	assert.True(selector.Matches(valid), selector.String())
	assert.True(selector.Matches(valid2))
	assert.False(selector.Matches(invalidMissing))
	assert.False(selector.Matches(invalid), selector.String())
}

func TestParseGroup(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"zoo":		"mar",
		"moo":		"lar",
		"thing":	"map",
	}
	invalid := Labels{
		"zoo":		"mar",
		"moo":		"something",
		"thing":	"map",
	}
	invalid2 := Labels{
		"zoo":		"mar",
		"moo":		"lar",
		"!thing":	"map",
	}
	selector, err := Parse("zoo=mar, moo=lar, thing")
	assert.Nil(err)
	assert.True(selector.Matches(valid))
	assert.False(selector.Matches(invalid))
	assert.False(selector.Matches(invalid2))

	complicated, err := Parse("zoo in (mar,lar,dar),moo,!thingy")
	assert.Nil(err)
	assert.NotNil(complicated)
	assert.True(complicated.Matches(valid))
}

func TestParseGroupComplicated(t *testing.T) {
	assert := assert.New(t)
	valid := Labels{
		"zoo":		"mar",
		"moo":		"lar",
		"thing":	"map",
	}
	complicated, err := Parse("zoo in (mar,lar,dar),moo,thing == map,!thingy")
	assert.Nil(err)
	assert.NotNil(complicated)
	assert.True(complicated.Matches(valid))
}

func TestParseDocsExample(t *testing.T) {
	assert := assert.New(t)
	sel, err := Parse("x in (foo,,baz),y,z notin ()")
	assert.Nil(err)
	assert.NotNil(sel)
}

func TestParseSubdomainKey(t *testing.T) {

	assert := assert.New(t)
	sel, err := Parse("example.com/failure-domain == primary")
	assert.Nil(err)
	assert.NotNil(sel)
	assert.Equal("example.com/failure-domain == primary", sel.String())
	assert.True(sel.Matches(map[string]string{
		"bar":				"foo",
		"example.com/failure-domain":	"primary",
		"foo":				"bar",
	}))
}

func TestParseEqualsOperators(t *testing.T) {
	assert := assert.New(t)

	selector, err := Parse("notin=in")
	assert.Nil(err)

	typed, isTyped := selector.(Equals)
	assert.True(isTyped)
	assert.Equal("notin", typed.Key)
	assert.Equal("in", typed.Value)
}

func TestParseValidate(t *testing.T) {
	assert := assert.New(t)

	_, err := Parse("zoo=bar")
	assert.Nil(err)

	_, err = Parse("_zoo=bar")
	assert.NotNil(err)

	_, err = Parse("_zoo=_bar")
	assert.NotNil(err)

	_, err = Parse("zoo=bar,foo=_mar")
	assert.NotNil(err)
}

func TestParseRegressionCSVSymbols(t *testing.T) {
	assert := assert.New(t)

	sel, err := Parse("foo in (bar-bar, baz.baz, buzz_buzz), moo=boo")
	assert.Nil(err, "regression is csv values can have '-', '.' and '_' in them")
	assert.NotEmpty(sel.String())
}

func TestParseRegressionCSVEarlyTermination(t *testing.T) {
	assert := assert.New(t)

	sel, err := Parse("foo in (bar,), moo=boo")
	assert.Nil(err)
	assert.NotEmpty(sel.String())
}

func TestParseRegressionCSVMultipleSpaces(t *testing.T) {
	assert := assert.New(t)

	sel, err := Parse("foo in (bar,  foo), moo=boo")
	assert.Nil(err)
	assert.NotEmpty(sel.String())
}

func TestParseRegressionIn(t *testing.T) {
	assert := assert.New(t)

	_, err := Parse("foo in bar, buzz)")
	assert.NotNil(err)
}

func TestParseMultiByte(t *testing.T) {
	assert := assert.New(t)

	selector, err := Parse("함=수,목=록")	// number=number, number=rock
	assert.Nil(err)
	assert.NotNil(selector)

	typed, isTyped := selector.(And)
	assert.True(isTyped)
	assert.Len(typed, 2)
}

func TestParseOptions(t *testing.T) {
	assert := assert.New(t)

	selQuery := "bar=foo@bar"
	labels := Labels{
		"foo":	"bar",
		"bar":	"foo@bar",
	}

	sel, err := Parse(selQuery)
	assert.NotNil(err)
	assert.Nil(sel)

	sel, err = Parse(selQuery, SkipValidation)
	assert.Nil(err)
	assert.NotNil(sel)

	assert.True(sel.Matches(labels))
}

func BenchmarkParse(b *testing.B) {
	valid := Labels{
		"zoo":		"mar",
		"moo":		"lar",
		"thing":	"map",
	}

	for i := 0; i < b.N; i++ {
		selector, err := Parse("zoo in (mar,lar,dar),moo,!thingy")
		if err != nil {
			b.Fail()
		}
		if !selector.Matches(valid) {
			b.Fail()
		}
	}
}

func TestParse_FuzzRegressions(t *testing.T) {
	assert := assert.New(t)

	var sel Selector
	var err error
	testBadStrings := []string{
		"!0!0",
		"0!=0,!",
	}
	for _, str := range testBadStrings {
		_, err = Parse(str)
		assert.NotNil(err, str, err)
	}

	testGoodStrings := []string{
		"0,!0",
		"0 in (0), !0",
	}
	for _, str := range testGoodStrings {
		sel, err = Parse(str)
		assert.Nil(err)
		_, err = Parse(sel.String())
		assert.Nil(err)
	}
}

func TestParse_InRegression(t *testing.T) {
	its := assert.New(t)

	good0 := Labels{
		"role":	"job",
		"team":	"internal-engineering",
	}
	good1 := Labels{
		"role":	"job-worker",
		"team":	"internal-engineering",
	}
	bad0 := Labels{
		"role": "not-job",
	}

	sel, err := Parse("role in (job, job-worker),team=internal-engineering")
	its.Nil(err)

	its.True(sel.Matches(good0))
	its.True(sel.Matches(good1))
	its.False(sel.Matches(bad0))
}
