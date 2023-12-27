/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
