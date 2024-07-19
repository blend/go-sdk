/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package profanity

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

var (
	_ Rule = (*NoGenericDecls)(nil)
)

// NoGenericDecls returns a profanity error if a generic function or type declaration exists
type NoGenericDecls struct {
	Enabled bool `yaml:"enabled"`
}

// Validate implements validation for the rule.
func (ngd NoGenericDecls) Validate() error {
	return nil
}

// Check implements Rule.
func (ngd NoGenericDecls) Check(filename string, contents []byte) RuleResult {
	if !ngd.Enabled {
		return RuleResult{OK: true}
	}
	if filepath.Ext(filename) != ".go" {
		return RuleResult{OK: true}
	}

	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, filename, contents, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return RuleResult{Err: err}
	}

	var result *RuleResult
	ast.Inspect(fileAst, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		if result != nil {
			return false
		}
		switch nt := n.(type) {
		case *ast.FuncType:
			if nt.TypeParams != nil {
				result = &RuleResult{
					File:    filename,
					Line:    fset.Position(nt.TypeParams.Pos()).Line,
					Message: "Type params present on function declaration",
				}
				return false
			}
		case *ast.TypeSpec:
			if nt.TypeParams != nil {
				result = &RuleResult{
					File:    filename,
					Line:    fset.Position(nt.TypeParams.Pos()).Line,
					Message: "Type params present on type declaration",
				}
				return false
			}
		case *ast.InterfaceType:
			if nt.Methods != nil {
				for _, i := range nt.Methods.List {
					if i == nil {
						continue
					}
					ast.Inspect(i, func(n ast.Node) bool {
						if n == nil {
							return false
						}
						if result != nil {
							return false
						}
						switch nt := n.(type) {
						case *ast.BinaryExpr:
							if nt.Op == token.OR {
								result = &RuleResult{
									File:    filename,
									Line:    fset.Position(nt.Pos()).Line,
									Message: "Union present in interface",
								}
								return false
							}
						case *ast.UnaryExpr:
							if nt.Op == token.TILDE {
								result = &RuleResult{
									File:    filename,
									Line:    fset.Position(nt.Pos()).Line,
									Message: "Underlying type operator present in interface",
								}
								return false
							}
						}
						return true
					})
				}
			}
			return false
		}
		return true
	})
	if result != nil {
		return *result
	}
	return RuleResult{OK: true}
}

// Strings implements fmt.Stringer.
func (ngd NoGenericDecls) String() string {
	return "generic declarations"
}
