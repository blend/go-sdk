package profanity

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// Calls returns a profanity error if a given function is called in a package.
func Calls(fns ...Call) RuleFunc {
	return func(filename string, contents []byte) RuleResult {
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
					for _, fn := range fns {
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
				case *ast.Ident:
					for _, fn := range fns {
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
}

func isIdent(expr ast.Expr, ident string) bool {
	if ident == "" {
		return true
	}
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

// Call is a package and function name pair.
//
// If package is empty string, it is assumed that the function
// is local to the calling package or a builtin.
type Call struct {
	Package string `yaml:"package"`
	Func    string `yaml:"func"`
}

// String implements fmt.Stringer
func (c Call) String() string {
	if c.Package != "" {
		return c.Package + "." + c.Func
	}
	return c.Func
}
