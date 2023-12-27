/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package copyright

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_KnownExtensions_templates(t *testing.T) {
	its := assert.New(t)

	for _, ext := range KnownExtensions {
		_, ok := DefaultExtensionNoticeTemplates[ext]
		its.True(ok, fmt.Sprintf("%s should have a known template", ext))
	}
}

func Test_KnownExtensions_includeFiles(t *testing.T) {
	its := assert.New(t)

	anyIncludeFiles := func(value string) bool {
		for _, include := range DefaultIncludeFiles {
			if value == include {
				return true
			}
		}
		return false
	}
	for _, ext := range KnownExtensions {
		ok := anyIncludeFiles("*" + ext)
		its.True(ok, fmt.Sprintf("%s should be in the included files list", ext))
	}
}

func Test_tsImportsTagMatch(t *testing.T) {
	its := assert.New(t)

	its.Matches(tsReferenceTagsExpr, goldenTsReferenceTags)
	its.Matches(tsReferenceTagsExpr, tsReferenceTags)
	its.NotMatches(tsReferenceTagsExpr, goldenTs)
}
