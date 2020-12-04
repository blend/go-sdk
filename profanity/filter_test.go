package profanity

import (
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Filter_IsZero(t *testing.T) {
	its := assert.New(t)

	its.True(Filter{}.IsZero())
	its.False(Filter{
		Include: []string{"foo", "bar"}, // any of these
	}.IsZero())
	its.False(Filter{
		Exclude: []string{"foo", "bar"}, // any of these
	}.IsZero())
	its.False(Filter{
		Include: []string{"foo", "bar"}, // any of these
		Exclude: []string{"foo", "bar"}, // any of these
	}.IsZero())
}

func Test_Filter_Match_Include(t *testing.T) {
	its := assert.New(t)

	f := Filter{
		Include: []string{"foo", "bar"}, // any of these
	}

	var includeMatch, excludeMatch string
	includeMatch, excludeMatch = f.Match("foo bar buzz", strings.Contains)
	its.Equal("foo", includeMatch)
	its.Equal("", excludeMatch)
}

func Test_Filter_Match_Include_Exclude(t *testing.T) {
	its := assert.New(t)

	f := Filter{
		Include: []string{"foo", "bar"},   // any of these
		Exclude: []string{"buzz", "wuzz"}, // but not these
	}

	var includeMatch, excludeMatch string
	includeMatch, excludeMatch = f.Match("foo bar buzz", strings.Contains)
	its.Equal("foo", includeMatch)
	its.Equal("buzz", excludeMatch)
}

func Test_Filter_Match_EmptyInput(t *testing.T) {
	its := assert.New(t)

	f := Filter{
		Include: []string{"foo", "bar"},   // any of these
		Exclude: []string{"buzz", "wuzz"}, // but not these
	}

	var includeMatch, excludeMatch string
	includeMatch, excludeMatch = f.Match("", strings.Contains)
	its.Equal("", includeMatch)
	its.Equal("", excludeMatch)
}

func Test_Filter_AllowMatch_IncludeExclude(t *testing.T) {
	its := assert.New(t)

	f := Filter{
		Include: []string{"foo", "bar"},   // any of these
		Exclude: []string{"buzz", "wuzz"}, // but not these
	}

	its.True(f.AllowMatch("test", ""))
	its.False(f.AllowMatch("test", "not-test"))
	its.False(f.AllowMatch("", "not-test"))
}

func Test_Filter_AllowMatch_Include(t *testing.T) {
	its := assert.New(t)

	f := Filter{
		Include: []string{"foo", "bar"}, // any of these
	}

	its.True(f.AllowMatch("test", ""))
	its.True(f.AllowMatch("test", "not-test"))
	its.False(f.AllowMatch("", "not-test"))
}

func Test_Filter_AllowMatch_Exclude(t *testing.T) {
	its := assert.New(t)

	f := Filter{
		Exclude: []string{"foo", "bar"}, // any of these
	}

	its.True(f.AllowMatch("test", ""))
	its.False(f.AllowMatch("test", "not-test"))
	its.False(f.AllowMatch("", "not-test"))
}
