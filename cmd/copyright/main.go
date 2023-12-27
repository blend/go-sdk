/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/copyright"
)

var (
	flagFallbackNoticeTemplate   string
	flagExtensionNoticeTemplates flagStrings
	flagNoticeBodyTemplate       string
	flagCompany                  string
	flagYear                     int
	flagLicense                  string

	flagRestrictions           string
	flagRestrictionsOpenSource bool
	flagRestrictionsInternal   bool

	flagVerify bool
	flagInject bool
	flagRemove bool

	flagExcludes        flagStrings
	flagExcludesFrom    flagStrings
	flagExcludeDefaults bool
	flagIncludeFiles    flagStrings

	flagExitFirst bool
	flagQuiet     bool
	flagVerbose   bool
	flagDebug     bool
	flagShowDiff  bool
)

func init() {
	flag.BoolVar(&flagQuiet, "quiet", false, "If all output should be suppressed")
	flag.BoolVar(&flagVerbose, "verbose", false, "If verbose output should be shown")
	flag.BoolVar(&flagDebug, "debug", false, "If debug output should be shown")
	flag.BoolVar(&flagShowDiff, "show-diff", false, "If the text diff in verification output should be shown")

	flag.BoolVar(&flagExitFirst, "exit-first", false, "If the program should exit on the first verification error")

	flag.StringVar(&flagNoticeBodyTemplate, "notice-body-template", copyright.DefaultNoticeBodyTemplate, "The notice body template; will try as a file path first, then used as a literal value. This is the template to inject into the filetype specific template for each file.")
	flag.StringVar(&flagFallbackNoticeTemplate, "fallback-notice-template", "", "The fallback notice template; will try as a file path first, then used as a literal value. This is the full notice (i.e. filetype specific) to use if there is no built-in notice template for the filetype.")

	flag.StringVar(&flagCompany, "company", "", "The company name to use in templates as {{ .Company }}")
	flag.IntVar(&flagYear, "year", time.Now().UTC().Year(), "The year to use in templates as {{ .Year }}")
	flag.StringVar(&flagLicense, "license", copyright.DefaultOpenSourceLicense, "The license to use in templates as {{ .License }}")
	flag.StringVar(&flagRestrictions, "restrictions", copyright.DefaultRestrictionsInternal, "The restriction template to compile and insert in the notice body template as {{ .Restrictions }}")
	flag.BoolVar(&flagRestrictionsOpenSource, "restrictions-open-source", false, fmt.Sprintf("The restrictions should be the open source defaults (i.e. %q)", copyright.DefaultRestrictionsOpenSource))
	flag.BoolVar(&flagRestrictionsInternal, "restrictions-internal", false, fmt.Sprintf("The restrictions should be the internal defaults (i.e. %q)", copyright.DefaultRestrictionsInternal))

	flag.BoolVar(&flagVerify, "verify", false, "If we should validate notices are present (exclusive with -inject and -remove) (this is the default)")
	flag.BoolVar(&flagInject, "inject", false, "If we should inject the notice (exclusive with -verify and -remove)")
	flag.BoolVar(&flagRemove, "remove", false, "If we should remove the notice (exclusive with -verify and -inject)")

	flag.Var(&flagExtensionNoticeTemplates, "ext", "Extension specific notice template overrides overrides; should be in the form -ext=js=js_template.txt, can be multiple")

	flag.BoolVar(&flagExcludeDefaults, "exclude-defaults", true, "If we should add the exclude defaults (e.g. node_modules etc.)")
	flag.Var(&flagExcludes, "exclude", "Files or directories to exclude via glob match, can be multiple")
	flag.Var(&flagExcludesFrom, "excludes-from", "A file to read for globs to exclude (e.g. .gitignore), can be multiple")
	flag.Var(&flagIncludeFiles, "include-file", "Files to include via glob match, can be multiple")

	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), `blend source code copyright management cli

> copyright [--inject|--verify|--remove] [ROOT(s)...]

Verify, inject or remove copyright notices from files in a given tree.

By default, this tool verifies that copyright notices are present with no flags provided.

Headers are treated exactly; do not edit the headers once they've been injected.

To verify headers:

	> copyright
	- OR -
	> copyright --verify
	- OR -
	> copyright --verify ./critical
	- OR -
	> copyright --verify ./critical/foo.py

To inject headers:

	> copyright --inject
	- OR -
	> copyright --inject ./critical/foo.py

	- NOTE: you can run "--inject" multiple times; it will only add the header if it is not present.

To remove headers:

	> copyright --remove

If you have an old version of the header in your files, and you want to migrate to an updated version:

	- Save the existing header to a file, "notice.txt", including any newlines between the notice and code
	- Remove existing notices:
		> copyright --remove -ext=py=notice.txt
	- Then inject the new notice:
		> copyright --inject --include-file="*.py"
	- You should now have the new notice in your files, and "--inject" will honor it

`,
		)
		oldUsage()
	}

	flag.Parse()
}

func main() {
	ctx := context.Background()

	if flagNoticeBodyTemplate == "" {
		fmt.Fprintln(os.Stderr, "--notice provided is an empty string; cannot continue")
		os.Exit(1)
	}

	var roots []string
	if args := flag.Args(); len(args) > 0 {
		roots = args[:]
	} else {
		roots = []string{"."}
	}

	if flagExcludeDefaults {
		flagExcludes = append(flagExcludes, flagStrings(copyright.DefaultExcludes)...)
	}
	if _, err := os.Stat(copyright.DefaultImplicitConfigFile); err == nil {
		fmt.Fprintf(os.Stdout, "Including implicit glob exclusion file: %q\n", copyright.DefaultImplicitConfigFile)
		flagExcludesFrom = append(flagExcludesFrom, copyright.DefaultImplicitConfigFile)
	}
	for _, excludesFrom := range flagExcludesFrom {
		excludes, err := readExcludesFile(excludesFrom)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		flagExcludes = append(flagExcludes, excludes...)
	}

	if len(flagIncludeFiles) == 0 {
		flagIncludeFiles = flagStrings(copyright.DefaultIncludeFiles)
	}

	var restrictions string
	if flagRestrictionsOpenSource {
		restrictions = copyright.DefaultRestrictionsOpenSource
	} else if flagRestrictionsInternal {
		restrictions = copyright.DefaultRestrictionsInternal
	} else {
		restrictions = flagRestrictions
	}

	extensionNoticeTemplates := copyright.DefaultExtensionNoticeTemplates
	for _, extValue := range flagExtensionNoticeTemplates {
		ext, noticeTemplate, err := parseExtensionNoticeBodyTemplate(extValue)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		extensionNoticeTemplates[ext] = noticeTemplate
	}

	engine := copyright.Copyright{
		Config: copyright.Config{
			FallbackNoticeTemplate:   tryReadFile(flagFallbackNoticeTemplate),
			NoticeBodyTemplate:       tryReadFile(flagNoticeBodyTemplate),
			Company:                  flagCompany,
			Restrictions:             restrictions,
			Year:                     flagYear,
			License:                  flagLicense,
			ExtensionNoticeTemplates: extensionNoticeTemplates,
			Excludes:                 flagExcludes,
			IncludeFiles:             flagIncludeFiles,
			ExitFirst:                &flagExitFirst,
			Quiet:                    &flagQuiet,
			Verbose:                  &flagVerbose,
			Debug:                    &flagDebug,
			ShowDiff:                 &flagShowDiff,
		},
	}

	var actions []func(context.Context, string) error
	var actionLabels []string

	if flagRemove {
		actions = append(actions, engine.Remove)
		actionLabels = append(actionLabels, "remove")
	}
	if flagInject {
		actions = append(actions, engine.Inject)
		actionLabels = append(actionLabels, "inject")
	}
	if flagVerify {
		actions = append(actions, engine.Verify)
		actionLabels = append(actionLabels, "verify")
	}
	if len(actions) == 0 {
		actions = append(actions, engine.Verify)
		actionLabels = append(actionLabels, "verify")
	}

	for index, action := range actions {
		didFail := false
		actionLabel := actionLabels[index]

		for _, root := range roots {
			maybeFail(ctx, action, root, &didFail)
		}

		if didFail {
			if !flagQuiet {
				fmt.Printf("copyright %s %s!\nuse `copyright --inject` to add missing notices\n", actionLabel, ansi.Red("failed"))
			}
			os.Exit(1)
		}
		if !flagQuiet {
			fmt.Printf("copyright %s %s!\n", actionLabel, ansi.Green("ok"))
		}
	}
}

type flagStrings []string

func (fs flagStrings) String() string {
	return strings.Join(fs, ", ")
}

func (fs *flagStrings) Set(flagValue string) error {
	if flagValue == "" {
		return fmt.Errorf("invalid flag value; is empty")
	}
	*fs = append(*fs, flagValue)
	return nil
}

func tryReadFile(path string) string {
	contents, err := os.ReadFile(path)
	if err != nil {
		return path
	}
	return strings.TrimSpace(string(contents))
}

func readExcludesFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var output []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		output = append(output, strings.TrimSpace(scanner.Text()))
	}
	return output, nil
}

func parseExtensionNoticeBodyTemplate(extensionNoticeBodyTemplate string) (extension, noticeBodyTemplate string, err error) {
	parts := strings.SplitN(extensionNoticeBodyTemplate, "=", 2)
	if len(parts) < 2 {
		err = fmt.Errorf("invalid `-ext` value; %s", extensionNoticeBodyTemplate)
		return
	}
	extension = parts[0]
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	noticeBodyTemplate = tryReadFile(parts[1])
	return
}

func maybeFail(ctx context.Context, action func(context.Context, string) error, root string, didFail *bool) {
	err := action(ctx, root)
	if err != nil {
		if err == copyright.ErrFailure {
			*didFail = true
			return
		}
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
