/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

import (
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

var (
	_ Rule = (*GoImports)(nil)
)

// GoImports returns a profanity error if a given file contains
// any of a list of imports based on a glob match.
type GoImports struct {
	GlobFilter `yaml:",inline"`
}

// Check implements Rule.
func (gi GoImports) Check(filename string, contents []byte) RuleResult {
	if filepath.Ext(filename) != ".go" {
		return RuleResult{OK: true}
	}

	fset := token.NewFileSet()

	ast, err := parser.ParseFile(fset, filename, contents, parser.ImportsOnly)
	if err != nil {
		return RuleResult{Err: err}
	}

	var includeGlob, excludeGlob string
	var fileImportPath string
	for _, fileImport := range ast.Imports {
		fileImportPath = strings.ReplaceAll(fileImport.Path.Value, "\"", "")
		if includeGlob, excludeGlob = gi.Match(fileImportPath); includeGlob != "" && excludeGlob == "" {
			return RuleResult{
				File:		filename,
				Line:		fset.Position(fileImport.Pos()).Line,
				Message:	fmt.Sprintf("go imports glob: \"%s\"", includeGlob),
			}
		}
	}
	return RuleResult{OK: true}
}

// String implements fmt.Stringer.
func (gi GoImports) String() string {
	return fmt.Sprintf("go imports %s", gi.GlobFilter.String())
}
