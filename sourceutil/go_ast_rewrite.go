/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sourceutil

import (
	"context"
	"go/ast"
)

// GoAstRewrite returns a go ast visitor with a given set of options.
func GoAstRewrite(opts ...GoAstRewriteOption) GoAstVisitor {
	var rewriteOpts GoAstRewriteOptions
	for _, opt := range opts {
		opt(&rewriteOpts)
	}
	return func(ctx context.Context, n ast.Node) bool {
		return rewriteOpts.Apply(ctx, n)
	}
}

// GoIsPackageCall returns a filter that determines if a function is a given sel.Fn.
//
// It will only evaluate for function calls that use a package selector
// that is, function calls that have a selector.
func GoIsPackageCall(pkg, fn string) GoAstRewriteOption {
	return func(opts *GoAstRewriteOptions) {
		opts.Filter = func(ctx context.Context, n ast.Node) (visit, recurse bool) {
			if nt, ok := n.(*ast.CallExpr); ok {
				if ft, ok := nt.Fun.(*ast.SelectorExpr); ok {
					if exprIsName(ft.X, pkg) && exprIsName(ft.Sel, fn) {
						return true, false // visit, do not recurse
					}
				}
				return false, false // do not visit, do not recurse
			}
			return false, true // do not visit, do recurse
		}
	}
}

// GoIsCall returns a filter that determines if a function is a given name.
//
// It will only evaluate for function calls that appear local to the
// current package, that is, function calls that do not have a selector.
func GoIsCall(fn string) GoAstRewriteOption {
	return func(opts *GoAstRewriteOptions) {
		opts.Filter = func(_ context.Context, n ast.Node) (visit, recurse bool) {
			if nt, ok := n.(*ast.CallExpr); ok {
				if ft, ok := nt.Fun.(*ast.Ident); ok {
					if exprIsName(ft, fn) {
						return true, false // visit, do not recurse
					}
				}
				return false, false // do not visit, do not recurse
			}
			return false, true // do not visit, do recurse
		}
	}
}

// GoRewritePackageCall changes a given function as filtered by a filter
// to a given call noted by sel.Fn.
func GoRewritePackageCall(sel, fn string) GoAstRewriteOption {
	return func(opts *GoAstRewriteOptions) {
		opts.NodeVisitor = func(_ context.Context, n ast.Node) {
			if nt, ok := n.(*ast.CallExpr); ok {
				if ft, ok := nt.Fun.(*ast.SelectorExpr); ok {
					exprSetName(ft.X, sel)
					exprSetName(ft.Sel, fn)
				}
			}
		}
	}
}

// GoRewriteCall changes a given function as filtered by a filter
// to a given call noted by Fn.
func GoRewriteCall(fn string) GoAstRewriteOption {
	return func(opts *GoAstRewriteOptions) {
		opts.NodeVisitor = func(_ context.Context, n ast.Node) {
			if nt, ok := n.(*ast.CallExpr); ok {
				if ft, ok := nt.Fun.(*ast.Ident); ok {
					exprSetName(ft, fn)
				}
			}
		}
	}
}

// GoAstVisitor mutates an ast node.
type GoAstVisitor func(context.Context, ast.Node) bool

// GoAstRewriteOption the ast rewrite options.
type GoAstRewriteOption func(*GoAstRewriteOptions)

// GoAstFilter is a delegate type that filters ast nodes for visiting.
type GoAstFilter func(context.Context, ast.Node) (visit, recurse bool)

// GoAstNodeVisitor mutates a given node.
type GoAstNodeVisitor func(context.Context, ast.Node)

// GoAstRewriteOptions breaks the mutator out into field specific mutators.
type GoAstRewriteOptions struct {
	Filter      GoAstFilter
	NodeVisitor GoAstNodeVisitor
}

// Apply applies the options to the ast node.
func (opts GoAstRewriteOptions) Apply(ctx context.Context, node ast.Node) bool {
	if opts.Filter != nil {
		visit, recurse := opts.Filter(ctx, node)
		if visit {
			if opts.NodeVisitor != nil {
				opts.NodeVisitor(ctx, node)
			}
		}
		return recurse
	}
	if opts.NodeVisitor != nil {
		opts.NodeVisitor(ctx, node)
	}
	return true // if no filter, always recurse
}

func exprIsName(expr ast.Expr, name string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == name
}

func exprSetName(expr ast.Expr, name string) {
	id, ok := expr.(*ast.Ident)
	if ok {
		id.Name = name
	}
}
