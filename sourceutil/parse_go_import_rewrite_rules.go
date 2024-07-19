/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sourceutil

import "strings"

// ParseGoImportRewriteRules parses go import rewrite rules as strings.
func ParseGoImportRewriteRules(rewriteRules []string) (output []GoImportVisitor) {
	for _, rewriteRule := range rewriteRules {
		pieces := strings.SplitN(rewriteRule, "=", 2)
		output = append(output, GoImportRewritePrefix(pieces[0], pieces[1]))
	}
	return
}
