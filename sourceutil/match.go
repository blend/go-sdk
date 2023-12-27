/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sourceutil

import (
	"fmt"
	"regexp"
)

// MatchRemove removes a line if it matches a given expression.
func MatchRemove(corpus []byte, expr string) []byte {
	compiledExpr := regexp.MustCompile(expr)
	return compiledExpr.ReplaceAll(corpus, nil)
}

// MatchInject injects a given value after the any instances of a given expression.
func MatchInject(corpus []byte, expr, inject string) []byte {
	compiledExpr := regexp.MustCompile(expr)
	return compiledExpr.ReplaceAll(corpus, []byte(fmt.Sprintf("$0\n%s\n", inject)))
}
