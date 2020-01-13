package profanity

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestContainsAny(t *testing.T) {
	assert := assert.New(t)

	ruleFunc := ContainsAny("foo", "bar")

	assert.NotNil(ok(ruleFunc("", []byte(`aaa foo`))))
	assert.NotNil(ok(ruleFunc("", []byte(`foo aaa`))))
	assert.NotNil(ok(ruleFunc("", []byte(`aaa foo aaa`))))

	assert.NotNil(ok(ruleFunc("", []byte(`aaa bar`))))
	assert.NotNil(ok(ruleFunc("", []byte(`bar aaa`))))
	assert.NotNil(ok(ruleFunc("", []byte(`aaa bar aaa`))))

	assert.Nil(ok(ruleFunc("", []byte(``))))
	assert.Nil(ok(ruleFunc("", []byte(`aaa`))))
}

func ok(res RuleResult) error {
	if !res.OK {
		return fmt.Errorf("not ok")
	}
	return nil
}

func TestContainsAnyReportsLineNumber(t *testing.T) {
	assert := assert.New(t)

	ruleFunc := ContainsAny("foo")

	file := `111
222
foo
444
555
`

	res := ruleFunc("", []byte(file))
	assert.Equal(3, res.Line)
}
