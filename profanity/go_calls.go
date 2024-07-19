/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package profanity

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/validate"
)

var (
	_ Rule = (*GoCalls)(nil)
)

// GoCalls returns a profanity error if a given function is called in a package.
type GoCalls []GoCall

// Validate implements validation for the rule.
func (gc GoCalls) Validate() error {
	for _, c := range gc {
		if c.Package == "" && c.Func == "" {
			return validate.Error(validate.ErrStringRequired, nil)
		}
	}
	return nil
}

// Check implements Rule.
func (gc GoCalls) Check(filename string, contents []byte) RuleResult {
	if filepath.Ext(filename) != ".go" {
		return RuleResult{OK: true}
	}

	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, filename, contents, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return RuleResult{Err: err}
	}

	var results []RuleResult
	ast.Inspect(fileAst, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch nt := n.(type) {
		case *ast.CallExpr:
			switch ft := nt.Fun.(type) {
			case *ast.SelectorExpr:
				for _, fn := range gc {
					if isIdent(ft.X, fn.Package) && isIdent(ft.Sel, fn.Func) {
						var message string
						if fn.Package != "" {
							message = fmt.Sprintf("go file includes function call: \"%s.%s\"", fn.Package, fn.Func)
						} else {
							message = fmt.Sprintf("go file includes function call: %q", fn.Func)
						}
						results = append(results, RuleResult{
							File:    filename,
							Line:    fset.Position(ft.Pos()).Line,
							Message: message,
						})
						return false
					}
				}
				return false
			case *ast.Ident: // check package local functions and built-ins
				for _, fn := range gc {
					if fn.Package == "" {
						if isIdent(ft, fn.Func) {
							results = append(results, RuleResult{
								File:    filename,
								Line:    fset.Position(ft.Pos()).Line,
								Message: fmt.Sprintf("go file includes function call: %q", fn.Func),
							})
							return false
						}
					}
				}
				return false
			}
		}
		return true
	})
	if len(results) > 0 {
		return results[0]
	}
	return RuleResult{OK: true}
}

// Strings implements fmt.Stringer.
func (gc GoCalls) String() string {
	var tokens []string
	for _, call := range gc {
		tokens = append(tokens, call.String())
	}
	return fmt.Sprintf("go calls: %s", strings.Join(tokens, " "))
}

// GoCall is a package and function name pair.
//
// `Package` is the package selector, typically the last path
// segment of the import (ex. "github.com/foo/bar" would be "bar")
//
// `Func` is the function name.
//
// If package is empty string, it is assumed that the function
// is local to the calling package or a builtin.
type GoCall struct {
	Package string `yaml:"package"`
	Func    string `yaml:"func"`
}

// String implements fmt.Stringer
func (gc GoCall) String() string {
	if gc.Package != "" {
		return gc.Package + "." + gc.Func
	}
	return gc.Func
}

func isIdent(expr ast.Expr, ident string) bool {
	if ident == "" {
		return true
	}
	id, ok := expr.(*ast.Ident)
	return ok && stringutil.Glob(id.Name, ident)
}
