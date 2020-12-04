package yaml_test

import (
	"strconv"
)

// MustUnquote unquotes a string, panicing if there is an issue.
func MustUnquote(str string) string {
	value, err := strconv.Unquote(str)
	if err != nil {
		panic(err)
	}
	return value
}
