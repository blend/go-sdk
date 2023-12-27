/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package codeowners

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// New creates a new copyright engine with a given set of config options.
func New(options ...Option) *Codeowners {
	var c Codeowners
	for _, option := range options {
		option(&c)
	}
	return &c
}

// Option mutates the Codeowners instance.
type Option func(*Codeowners)

// Codeowners holds the engine that generates and validates codeowners files.
type Codeowners struct {
	// Config holds the configuration opitons.
	Config

	// Stdout is the writer for Verbose and Debug output.
	// If it is unset, `os.Stdout` will be used.
	Stdout	io.Writer
	// Stderr is the writer for Error output.
	// If it is unset, `os.Stderr` will be used.
	Stderr	io.Writer
}

// GenerateFile generates the file as nominated by the config path.
func (c Codeowners) GenerateFile(ctx context.Context, root string) error {
	f, err := os.Create(c.PathOrDefault())
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	c.Debugf("codeowners path: %s", c.PathOrDefault())
	return c.Generate(ctx, root, f)
}

// Generate generates a codeowner file.
func (c Codeowners) Generate(ctx context.Context, root string, output io.Writer) error {
	var codeowners File
	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// skip common bogus dirs
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), "_") {
				return filepath.SkipDir
			}
			if info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			if strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			if strings.HasPrefix(path, "vendor/") {
				return filepath.SkipDir
			}
			return nil
		}

		// handle go files specially
		if strings.HasSuffix(info.Name(), ".go") {
			owners, parseErr := ParseGoComments(root, path, OwnersGoCommentPrefix)
			if parseErr != nil {
				return parseErr
			}
			if owners != nil {
				codeowners = append(codeowners, *owners)
			}
			return nil
		}

		// handle the owners file specially
		if info.Name() == OwnersFile {
			parsed, parseErr := ParseSource(root, path)
			if parseErr != nil {
				return parseErr
			}
			if parsed != nil {
				codeowners = append(codeowners, *parsed)
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	_, err = codeowners.WriteTo(output)
	return err
}

// ValidateFile validates the file as configured in the config field.
func (c Codeowners) ValidateFile(ctx context.Context) error {
	f, err := os.Open(c.PathOrDefault())
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	c.Debugf("codeowners path: %s", c.PathOrDefault())
	return c.Validate(ctx, f)
}

// Validate validates a given codeowners file.
func (c Codeowners) Validate(ctx context.Context, input io.Reader) error {
	if c.GithubToken == "" {
		return fmt.Errorf("codeowners cannot validate; github token is empty")
	}

	codeownersFile, err := Read(input)
	if err != nil {
		return err
	}
	ghc := GithubClient{
		Addr:	c.Config.GithubURLOrDefault(),
		Token:	c.Config.GithubToken,
	}
	for _, source := range codeownersFile {
		for _, path := range source.Paths {
			// test that the path
			pathGlob := strings.TrimSuffix(path.PathGlob, "**")
			pathGlob = strings.TrimSuffix(pathGlob, "*")
			pathGlob = strings.TrimPrefix(pathGlob, "/")
			pathGlob = filepath.Join("./", pathGlob)
			if _, err := os.Stat(pathGlob); err != nil {
				return fmt.Errorf("codeowners path glob doesn't exist: %q", pathGlob)
			}

			// test that the owner(s) exist in github
			for _, owner := range path.Owners {
				c.Verbosef("codeowners source: %s; checking if owner exists: %s", source.Source, owner)
				if strings.Contains(owner, "/") {
					if err := ghc.TeamExists(ctx, owner); err != nil {
						return fmt.Errorf("github team not found: %q", owner)
					}
				} else {
					if err := ghc.UserExists(ctx, owner); err != nil {
						return fmt.Errorf("github user not found: %q", owner)
					}
				}
			}
		}
	}
	return nil
}

// GetStdout returns standard out.
func (c Codeowners) GetStdout() io.Writer {
	if c.QuietOrDefault() {
		return io.Discard
	}
	if c.Stdout != nil {
		return c.Stdout
	}
	return os.Stdout
}

// GetStderr returns standard error.
func (c Codeowners) GetStderr() io.Writer {
	if c.QuietOrDefault() {
		return io.Discard
	}
	if c.Stderr != nil {
		return c.Stderr
	}
	return os.Stderr
}

// Verbosef writes to stdout if the `Verbose` flag is true.
func (c Codeowners) Verbosef(format string, args ...interface{}) {
	if !c.VerboseOrDefault() {
		return
	}
	fmt.Fprintf(c.GetStdout(), format+"\n", args...)
}

// Debugf writes to stdout if the `Debug` flag is true.
func (c Codeowners) Debugf(format string, args ...interface{}) {
	if !c.DebugOrDefault() {
		return
	}
	fmt.Fprintf(c.GetStdout(), format+"\n", args...)
}
