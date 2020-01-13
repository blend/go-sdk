package selector

import (
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestParserIsWhitespace(t *testing.T) {
	assert := assert.New(t)

	l := &Parser{}
	assert.True(l.isWhitespace(' '))
	assert.True(l.isWhitespace('\n'))
	assert.True(l.isWhitespace('\r'))
	assert.True(l.isWhitespace('\t'))

	assert.False(l.isWhitespace('a'))
	assert.False(l.isWhitespace('z'))
	assert.False(l.isWhitespace('A'))
	assert.False(l.isWhitespace('Z'))
	assert.False(l.isWhitespace('1'))
	assert.False(l.isWhitespace('-'))
}

func TestParserIsAlpha(t *testing.T) {
	assert := assert.New(t)

	l := &Parser{}
	assert.True(l.isAlpha('a'))
	assert.True(l.isAlpha('z'))
	assert.True(l.isAlpha('A'))
	assert.True(l.isAlpha('Z'))
	assert.True(l.isAlpha('1'))

	assert.False(l.isAlpha('-'))
	assert.False(l.isAlpha(' '))
	assert.False(l.isAlpha('\n'))
	assert.False(l.isAlpha('\r'))
	assert.False(l.isAlpha('\t'))
}

func TestParserSkipWhitespace(t *testing.T) {
	assert := assert.New(t)

	l := &Parser{s: "foo    != bar    ", pos: 3}
	assert.Equal(" ", string(l.current()))
	l.skipWhiteSpace()
	assert.Equal(7, l.pos)
	assert.Equal("!", string(l.current()))
	l.pos = 14
	assert.Equal(" ", string(l.current()))
	l.skipWhiteSpace()
	assert.Equal(len(l.s), l.pos)
}

func TestParserReadWord(t *testing.T) {
	assert := assert.New(t)

	l := &Parser{s: "foo != bar"}
	assert.Equal("foo", l.readWord())
	assert.Equal(" ", string(l.current()))

	l = &Parser{s: "foo,"}
	assert.Equal("foo", l.readWord())
	assert.Equal(",", string(l.current()))

	l = &Parser{s: "foo"}
	assert.Equal("foo", l.readWord())
	assert.True(l.done())
}

func TestParserReadOp(t *testing.T) {
	assert := assert.New(t)

	l := &Parser{s: "!= bar"}
	op, err := l.readOp()
	assert.Nil(err)
	assert.Equal("!=", op)
	assert.Equal(" ", string(l.current()))

	l = &Parser{s: "!=bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("!=", op)
	assert.Equal("b", string(l.current()))

	l = &Parser{s: "!=bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("!=", op)
	assert.Equal("b", string(l.current()))

	l = &Parser{s: "!="}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("!=", op)
	assert.True(l.done())

	l = &Parser{s: "= bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("=", op)
	assert.Equal(" ", string(l.current()))

	l = &Parser{s: "=bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("=", op)
	assert.Equal("b", string(l.current()))

	l = &Parser{s: "== bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("==", op)
	assert.Equal(" ", string(l.current()))

	l = &Parser{s: "==bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("==", op)
	assert.Equal("b", string(l.current()))

	l = &Parser{s: "in (foo)"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("in", op)
	assert.Equal(" ", string(l.current()))

	l = &Parser{s: "in(foo)"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("in", op)
	assert.Equal("(", string(l.current()))

	l = &Parser{s: "notin (foo)"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("notin", op)
	assert.Equal(" ", string(l.current()))

	l = &Parser{s: "notin(foo)"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("notin", op)
	assert.Equal("(", string(l.current()))
}

func TestParserReadCSV(t *testing.T) {
	assert := assert.New(t)

	l := &Parser{s: "(bar, baz, biz)"}
	words, err := l.readCSV()
	assert.Nil(err)
	assert.Len(words, 3, strings.Join(words, ","))
	assert.Equal("bar", words[0])
	assert.Equal("baz", words[1])
	assert.Equal("biz", words[2])
	assert.True(l.done())

	l = &Parser{s: "(bar,baz,biz)"}
	words, err = l.readCSV()
	assert.Nil(err)
	assert.Len(words, 3, strings.Join(words, ","))
	assert.Equal("bar", words[0])
	assert.Equal("baz", words[1])
	assert.Equal("biz", words[2])
	assert.True(l.done())

	l = &Parser{s: "(bar, buzz, baz"}
	words, err = l.readCSV()
	assert.NotNil(err)

	l = &Parser{s: "()"}
	words, err = l.readCSV()
	assert.Nil(err)
	assert.Empty(words)
	assert.True(l.done())

	l = &Parser{s: "(), thing=after"}
	words, err = l.readCSV()
	assert.Nil(err)
	assert.Empty(words)
	assert.Equal(",", string(l.current()))

	l = &Parser{s: "(foo, bar), buzz=light"}
	words, err = l.readCSV()
	assert.Nil(err)
	assert.Len(words, 2)
	assert.Equal("foo", words[0])
	assert.Equal("bar", words[1])
	assert.Equal(",", string(l.current()))

	l = &Parser{s: "(test, space are bad)"}
	words, err = l.readCSV()
	assert.NotNil(err)
}

func TestParserHasKey(t *testing.T) {
	assert := assert.New(t)
	l := &Parser{s: "foo"}
	valid, err := l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(HasKey)
	assert.True(isTyped)
	assert.Equal("foo", string(typed))
}

func TestParserNotHasKey(t *testing.T) {
	assert := assert.New(t)
	l := &Parser{s: "!foo"}
	valid, err := l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(NotHasKey)
	assert.True(isTyped)
	assert.Equal("foo", string(typed))
}

func TestParserEquals(t *testing.T) {
	assert := assert.New(t)

	l := &Parser{s: "foo = bar"}
	valid, err := l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(Equals)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Equal("bar", typed.Value)

	l = &Parser{s: "foo=bar"}
	valid, err = l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped = valid.(Equals)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Equal("bar", typed.Value)
}

func TestParserDoubleEquals(t *testing.T) {
	assert := assert.New(t)
	l := &Parser{s: "foo == bar"}
	valid, err := l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(Equals)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Equal("bar", typed.Value)
}

func TestParserNotEquals(t *testing.T) {
	assert := assert.New(t)
	l := &Parser{s: "foo != bar"}
	valid, err := l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(NotEquals)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Equal("bar", typed.Value)
}

func TestParserIn(t *testing.T) {
	assert := assert.New(t)
	l := &Parser{s: "foo in (bar, baz)"}
	valid, err := l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(In)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Len(typed.Values, 2)
	assert.Equal("bar", typed.Values[0])
	assert.Equal("baz", typed.Values[1])
}

func TestParserLex(t *testing.T) {
	assert := assert.New(t)
	l := &Parser{s: ""}
	_, err := l.Parse()
	assert.Nil(err)
}
