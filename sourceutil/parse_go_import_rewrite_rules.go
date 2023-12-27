/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
