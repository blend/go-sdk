/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sourceutil

import (
	"fmt"
)

// GoImportRewritePrefix is a helper that returns a visitor that
// takes a source path and rewrites it as the destination path
// preserving any path segments after the source if it matches the sourcePrefix.
func GoImportRewritePrefix(sourcePrefix, destinationPrefix string) GoImportVisitor {
	return GoImportRewrite(
		OptGoImportPathRewrite(
			fmt.Sprintf("^%s(.*)", sourcePrefix),
			fmt.Sprintf("%s$1", destinationPrefix),
		),
	)
}
