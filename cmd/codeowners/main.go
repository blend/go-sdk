/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/codeowners"
)

var (
	flagPath        string
	flagGithubURL   string
	flagGithubToken string

	flagValidate bool
	flagGenerate bool

	flagQuiet   bool
	flagVerbose bool
	flagDebug   bool
)

func init() {
	flag.StringVar(&flagPath, "path", codeowners.DefaultPath, "The codeowners file path")
	flag.StringVar(&flagGithubURL, "github-url", codeowners.DefaultGithubURL, "The github api url")
	flag.StringVar(&flagGithubToken, "github-token", os.Getenv(codeowners.DefaultGithubTokenEnvVar), "The github api token")

	flag.BoolVar(&flagQuiet, "quiet", false, "If all output should be suppressed")
	flag.BoolVar(&flagVerbose, "verbose", false, "If verbose output should be shown")
	flag.BoolVar(&flagDebug, "debug", false, "If debug output should be shown")

	flag.BoolVar(&flagValidate, "validate", false, "If we should validate the codeowners file (exclusive with -generate) (this is the default)")
	flag.BoolVar(&flagGenerate, "generate", false, "If we should generate the codeowners file (exclusive with -validate)")

	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), `github codeowners management cli

Verify or generate github codeowners files.

By default, this tool verifies that the codeowners file is present, and that it
contains valid user and team references. 

To generate the codeowners file:

	> codeowners --generate

This will scan the file tree for '.codeowners' file that relatively nominate
owners for a given directory and child directories.

You can also annotate files with the //github:codeowner source comment.
`,
		)
		oldUsage()
	}

	flag.Parse()
}

func main() {
	ctx := context.Background()

	engine := codeowners.Codeowners{
		Config: codeowners.Config{
			Path:        flagPath,
			GithubURL:   flagGithubURL,
			GithubToken: flagGithubToken,
			Quiet:       &flagQuiet,
			Verbose:     &flagVerbose,
			Debug:       &flagDebug,
		},
	}

	var actionLabel string
	var err error
	if flagGenerate {
		actionLabel = "generate"
		var root string
		if args := flag.Args(); len(args) > 0 {
			root = args[0]
		} else {
			root = "."
		}
		err = engine.GenerateFile(ctx, root)
	} else if flagValidate {
		actionLabel = "validate"
		err = engine.ValidateFile(ctx)
	} else { // the default
		actionLabel = "validate"
		err = engine.ValidateFile(ctx)
	}

	if err != nil {
		if !flagQuiet {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			fmt.Printf("codeowners %s %s!\n", actionLabel, ansi.Red("failed"))
		}
		os.Exit(1)
	}
	if !flagQuiet {
		fmt.Printf("codeowners %s %s!\n", actionLabel, ansi.Green("ok"))
	}
}
