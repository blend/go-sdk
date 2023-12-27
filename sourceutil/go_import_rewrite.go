/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sourceutil

import (
	"context"
	"go/ast"
	"regexp"
	"strings"
)

// GoImportRewrite visits and optionally mutates imports for go files.
func GoImportRewrite(opts ...GoImportRewriteOption) GoImportVisitor {
	var rewriteOpts GoImportRewriteOptions
	for _, opt := range opts {
		opt(&rewriteOpts)
	}
	return func(ctx context.Context, importSpec *ast.ImportSpec) error {
		return rewriteOpts.Apply(ctx, importSpec)
	}
}

// OptGoImportPathMatches returns a rewrite filter that returns if an import path matches a given expression.
func OptGoImportPathMatches(expr string) GoImportRewriteOption {
	return func(opts *GoImportRewriteOptions) {
		opts.Filter = func(ctx context.Context, importSpec *ast.ImportSpec) (bool, error) {
			compiledExpr, err := regexp.Compile(expr)
			if err != nil {
				return false, err
			}
			return compiledExpr.MatchString(RemoveQuotes(importSpec.Path.Value)), nil
		}
	}
}

// OptGoImportNameMatches returns a rewrite filter that returns if an import name matches a given expression.
func OptGoImportNameMatches(expr string) GoImportRewriteOption {
	return func(opts *GoImportRewriteOptions) {
		opts.Filter = func(ctx context.Context, importSpec *ast.ImportSpec) (bool, error) {
			compiledExpr, err := regexp.Compile(expr)
			if err != nil {
				return false, err
			}
			return compiledExpr.MatchString(importSpec.Name.Name), nil
		}
	}
}

// OptGoImportAddName adds a name if one is not already specified.
func OptGoImportAddName(name string) GoImportRewriteOption {
	return func(opts *GoImportRewriteOptions) {
		opts.NameVisitor = func(ctx context.Context, nameNode *ast.Ident) error {
			if nameNode.Name == "" {
				nameNode.Name = name
			}
			return nil
		}
	}
}

// OptGoImportSetAlias sets the import alias to a given value.
//
// Setting to "" will remove the alias.
func OptGoImportSetAlias(name string) GoImportRewriteOption {
	return func(opts *GoImportRewriteOptions) {
		opts.NameVisitor = func(ctx context.Context, nameNode *ast.Ident) error {
			nameNode.Name = name
			return nil
		}
	}
}

// OptGoImportSetPath sets an import path to a given value.
func OptGoImportSetPath(path string) GoImportRewriteOption {
	return func(opts *GoImportRewriteOptions) {
		opts.PathVisitor = func(ctx context.Context, pathNode *ast.BasicLit) error {
			pathNode.Value = path
			if !strings.HasPrefix(pathNode.Value, "\"") {
				pathNode.Value = "\"" + pathNode.Value
			}
			if !strings.HasSuffix(pathNode.Value, "\"") {
				pathNode.Value = pathNode.Value + "\""
			}
			return nil
		}
	}
}

// OptGoImportPathRewrite returns a path filter and rewrite expression.
func OptGoImportPathRewrite(matchExpr, outputExpr string) GoImportRewriteOption {
	return func(opts *GoImportRewriteOptions) {
		compiledMatch, compileErr := regexp.Compile(matchExpr)
		opts.Filter = func(ctx context.Context, importSpec *ast.ImportSpec) (output bool, err error) {
			if compileErr != nil {
				err = compileErr
				return
			}
			importPath := RemoveQuotes(importSpec.Path.Value)
			output = compiledMatch.MatchString(importPath)
			return
		}
		opts.PathVisitor = func(ctx context.Context, path *ast.BasicLit) error {
			output := []byte{}
			for _, submatches := range compiledMatch.FindAllStringSubmatchIndex(RemoveQuotes(path.Value), -1) {
				output = compiledMatch.ExpandString(output, outputExpr, RemoveQuotes(path.Value), submatches)
			}
			path.Value = string(output)
			if !strings.HasPrefix(path.Value, "\"") {
				path.Value = "\"" + path.Value
			}
			if !strings.HasSuffix(path.Value, "\"") {
				path.Value = path.Value + "\""
			}
			return nil
		}
	}
}

// GoImportVisitor mutates an ast import.
type GoImportVisitor func(context.Context, *ast.ImportSpec) error

// GoImportRewriteOption mutates the import rewrite options
type GoImportRewriteOption func(*GoImportRewriteOptions)

// GoImportRewriteOptions breaks the mutator out into field specific mutators.
type GoImportRewriteOptions struct {
	Filter		func(context.Context, *ast.ImportSpec) (bool, error)
	CommentVisitor	func(context.Context, *ast.CommentGroup) error
	DocVisitor	func(context.Context, *ast.CommentGroup) error
	NameVisitor	func(context.Context, *ast.Ident) error
	PathVisitor	func(context.Context, *ast.BasicLit) error
}

// Apply applies the options to the import.
func (opts GoImportRewriteOptions) Apply(ctx context.Context, importSpec *ast.ImportSpec) error {
	if opts.Filter != nil {
		if ok, err := opts.Filter(ctx, importSpec); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	if opts.CommentVisitor != nil {
		if importSpec.Comment == nil {
			importSpec.Comment = &ast.CommentGroup{}
		}
		if err := opts.CommentVisitor(ctx, importSpec.Comment); err != nil {
			return err
		}
	}
	if opts.DocVisitor != nil {
		if importSpec.Doc == nil {
			importSpec.Doc = &ast.CommentGroup{}
		}
		if err := opts.DocVisitor(ctx, importSpec.Doc); err != nil {
			return err
		}
	}
	if opts.NameVisitor != nil {
		if importSpec.Name == nil {
			importSpec.Name = &ast.Ident{}
		}
		if err := opts.NameVisitor(ctx, importSpec.Name); err != nil {
			return err
		}
	}
	if opts.PathVisitor != nil {
		if err := opts.PathVisitor(ctx, importSpec.Path); err != nil {
			return err
		}
	}
	return nil
}
