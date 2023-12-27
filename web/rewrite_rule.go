/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import "regexp"

// RewriteAction is an action for a rewrite rule.
type RewriteAction func(filePath string, matchedPieces ...string) string

// RewriteRule is a rule for re-writing incoming static urls.
type RewriteRule struct {
	MatchExpression	string
	expr		*regexp.Regexp
	Action		RewriteAction
}

// Apply runs the filter, returning a bool if it matched, and the resulting path.
func (rr RewriteRule) Apply(filePath string) (bool, string) {
	if rr.expr.MatchString(filePath) {
		pieces := extractSubMatches(rr.expr, filePath)
		return true, rr.Action(filePath, pieces...)
	}

	return false, filePath
}

// ExtractSubMatches returns sub matches for an expr because go's regexp library is weird.
func extractSubMatches(re *regexp.Regexp, corpus string) []string {
	allResults := re.FindAllStringSubmatch(corpus, -1)
	results := []string{}
	for _, resultSet := range allResults {
		results = append(results, resultSet...)
	}
	return results
}
