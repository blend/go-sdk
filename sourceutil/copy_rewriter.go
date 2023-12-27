/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sourceutil

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/blend/go-sdk/stringutil"
)

// CopyRewriter copies a source to a destination, and applies rewrite rules to the file(s) it copies.
type CopyRewriter struct {
	Source			string
	Destination		string
	SkipGlobs		[]string
	GoImportVisitors	[]GoImportVisitor
	GoAstVistiors		[]GoAstVisitor
	StringSubstitutions	[]StringSubstitution
	DryRun			bool
	RemoveDestination	bool
	KeepTemporary		bool

	Quiet	*bool
	Verbose	*bool
	Debug	*bool

	Stdout	io.Writer
	Stderr	io.Writer
}

// Execute is the command body.
func (cr CopyRewriter) Execute(ctx context.Context) error {
	if _, err := os.Stat(cr.Source); err != nil {
		return fmt.Errorf("source not found at %s", cr.Source)
	}
	tempDir, err := os.MkdirTemp("", "repoctl")
	if err != nil {
		return err
	}
	if !cr.KeepTemporary {
		defer func() {
			if _, err = os.Stat(tempDir); err == nil {
				cr.Verbosef("cleaning up temp dir %s", tempDir)
				os.RemoveAll(tempDir)
			}
		}()
	}

	// walk files
	err = filepath.Walk(cr.Source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		base := strings.TrimPrefix(strings.TrimPrefix(path, cr.Source), "/")
		destination := filepath.Join(tempDir, base)

		for _, skipGlob := range cr.SkipGlobs {
			if stringutil.Glob(base, skipGlob) {
				if info.IsDir() {
					cr.Verbosef("%s: skipping dir", base)
					return filepath.SkipDir
				}
				cr.Verbosef("%s: skipping", base)
				return nil
			}
		}

		if info.IsDir() {
			if _, err := os.Stat(destination); err != nil {
				cr.Verbosef("%s", base)
				if !cr.DryRun {
					cr.Debugf("%s: creating %s", base, destination)
					if err = os.MkdirAll(destination, DefaultDirPerms); err != nil {
						return err
					}
				} else {
					cr.Debugf("%s: dry-run; creating dir %s", base, destination)
				}
			}
			return nil
		}

		cr.Verbosef("%s", base)
		if filepath.Ext(path) == ".go" {
			if err := cr.copyGoSourceFile(ctx, destination, path); err != nil {
				return err
			}
		} else {
			if !cr.DryRun {
				if err := Copy(ctx, destination, path); err != nil {
					return err
				}
			}
		}
		return nil
	})

	if !cr.DryRun {
		if cr.RemoveDestination {
			cr.Verbosef("removing destination dir %s", cr.Destination)
			if err := os.RemoveAll(cr.Destination); err != nil {
				return err
			}
		}
		cr.Verbosef("recursively copying %s to %s", tempDir, cr.Destination)
		if err := CopyAll(cr.Destination, tempDir); err != nil {
			return err
		}
	} else {
		cr.Verbosef("%s", "dry-run; skipping final copy")
	}
	return nil
}

// copyGoSourceFile rewrites the imports for a golang file at a given path
func (cr CopyRewriter) copyGoSourceFile(ctx context.Context, destinationPath, sourcePath string) error {
	contents, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	var writer io.WriteCloser
	if cr.DryRun {
		writer = nopWriteCloser{io.Discard}
	} else {
		writer, err = os.Create(destinationPath)
		if err != nil {
			return err
		}
		defer writer.Close()
	}
	if err = cr.rewriteGoAst(ctx, sourcePath, contents, writer); err != nil {
		return err
	}
	return cr.rewriteContents(ctx, destinationPath)
}

func (cr CopyRewriter) rewriteGoAst(ctx context.Context, sourcePath string, contents []byte, writer io.Writer) error {
	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, sourcePath, contents, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return err
	}

	for importIndex := range fileAst.Imports {	// foreach file import
		cr.Debugf("processing import %s", fileAst.Imports[importIndex].Path.Value)
		for _, rewriteRule := range cr.GoImportVisitors {	// foreach import rule
			if err := rewriteRule(ctx, fileAst.Imports[importIndex]); err != nil {
				return err
			}
		}
	}
	for _, rewrite := range cr.GoAstVistiors {
		ast.Inspect(fileAst, func(n ast.Node) bool {
			if n == nil {
				return false
			}
			return rewrite(ctx, n)
		})
	}
	return printer.Fprint(writer, fset, fileAst)
}

func (cr CopyRewriter) rewriteContents(ctx context.Context, sourcePath string) error {
	if len(cr.StringSubstitutions) == 0 {
		return nil
	}

	stat, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	contents, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	var output string
	var ok bool
	for _, rule := range cr.StringSubstitutions {
		output, ok = rule(ctx, string(contents))
		if ok {
			contents = []byte(output)
		}
	}
	if cr.DryRun {
		cr.Debugf("dry-run; skipping rewriting file %s", sourcePath)
		return nil
	}
	cr.Debugf("rewriting file %s", sourcePath)
	return os.WriteFile(sourcePath, contents, stat.Mode())
}

// QuietOrDefault returns a value or a default.
func (cr CopyRewriter) QuietOrDefault() bool {
	if cr.Quiet != nil {
		return *cr.Quiet
	}
	return false
}

// VerboseOrDefault returns a value or a default.
func (cr CopyRewriter) VerboseOrDefault() bool {
	if cr.Verbose != nil {
		return *cr.Verbose
	}
	return false
}

// DebugOrDefault returns a value or a default.
func (cr CopyRewriter) DebugOrDefault() bool {
	if cr.Debug != nil {
		return *cr.Debug
	}
	return false
}

// GetStdout returns standard out.
func (cr CopyRewriter) GetStdout() io.Writer {
	if cr.QuietOrDefault() {
		return io.Discard
	}
	if cr.Stdout != nil {
		return cr.Stdout
	}
	return os.Stdout
}

// GetStderr returns standard error.
func (cr CopyRewriter) GetStderr() io.Writer {
	if cr.QuietOrDefault() {
		return io.Discard
	}
	if cr.Stderr != nil {
		return cr.Stderr
	}
	return os.Stderr
}

// Verbosef writes to stdout if the `Verbose` flag is true.
func (cr CopyRewriter) Verbosef(format string, args ...interface{}) {
	if !cr.VerboseOrDefault() {
		return
	}
	fmt.Fprintf(cr.GetStdout(), format+"\n", args...)
}

// Debugf writes to stdout if the `Debug` flag is true.
func (cr CopyRewriter) Debugf(format string, args ...interface{}) {
	if !cr.DebugOrDefault() {
		return
	}
	fmt.Fprintf(cr.GetStdout(), format+"\n", args...)
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error	{ return nil }
