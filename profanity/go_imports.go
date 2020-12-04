package profanity

import (
	"fmt"
	"go/parser"
	"go/token"
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
				File:    filename,
				Line:    fset.Position(fileImport.Pos()).Line,
				Message: fmt.Sprintf("go imports glob: \"%s\"", includeGlob),
			}
		}
	}
	return RuleResult{OK: true}
}

// String implements fmt.Stringer.
func (gi GoImports) String() string {
	return fmt.Sprintf("go imports %s", gi.GlobFilter.String())
}
