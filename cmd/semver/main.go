package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blend/go-sdk/semver"
)

// linker metadata block
// this block must be present
// it is used by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func usage() {
	fmt.Fprint(os.Stdout, "version validates and manage versions from a given file\n")
	fmt.Fprint(os.Stdout, "\ncommands:\n")
	fmt.Fprint(os.Stdout, "\tvalidate\t\t\t\tvalidate a given version file\n")
	fmt.Fprint(os.Stdout, "\tincrement major|minor|patch\t\tincrement a given version segment\n")
	fmt.Fprint(os.Stdout, "\tsatisfies <version constraint>\t\tverify a version satisfies a constraint\n")
	fmt.Fprint(os.Stdout, "\nusage:\n")
	fmt.Fprint(os.Stdout, "\tversion [command] [args] -f [filename]\n")
}

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	command, args := os.Args[1], os.Args[2:]

	switch command {
	case "increment":
		increment(args)
		os.Exit(0)
	case "satisfies":
		satisfies(args)
		os.Exit(0)
	case "validate":
		validate(args)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "invalid command: %s\n", command)
		os.Exit(1)
	}
}

func readContents(path string) (contents []byte, err error) {
	if strings.TrimSpace(path) == "-" {
		contents, err = ioutil.ReadAll(os.Stdin)
	} else {
		contents, err = ioutil.ReadFile(path)
	}
	return
}

func increment(args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "must supply a semver segment and a file\n")
		os.Exit(1)
	}

	segment := args[0]
	filepath := args[1]

	contents, err := readContents(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	version, err := semver.NewVersion(strings.TrimSpace(string(contents)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	switch strings.ToLower(segment) {
	case "patch":
		version.BumpPatch()
	case "minor":
		version.BumpMinor()
	case "major":
		version.BumpMajor()
	default:
		fmt.Fprintf(os.Stderr, "invalid segment: %s\n", segment)
		fmt.Fprintf(os.Stderr, "should be one of: 'major', 'minor', and 'patch'\n")
		os.Exit(1)
	}

	fmt.Printf("%v\n", version)
}

func satisfies(args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "must supply a version constraint and a file\n")
		os.Exit(1)
	}

	constraintValue := args[0]
	filepath := args[1]

	contents, err := readContents(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	version, err := semver.NewVersion(strings.TrimSpace(string(contents)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	constraint, err := semver.NewConstraint(constraintValue)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	if !constraint.Check(version) {
		fmt.Fprintf(os.Stderr, "`%s` does not satisfy `%s`\n", constraint.String(), version.String())
		os.Exit(1)
	}
}

func validate(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "must supply a file\n")
		os.Exit(1)
	}
	filepath := args[0]
	contents, err := readContents(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	_, err = semver.NewVersion(strings.TrimSpace(string(contents)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
