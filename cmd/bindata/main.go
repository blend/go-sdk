package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blend/go-sdk/bindata"
	"github.com/blend/go-sdk/ex"
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Short: "bindata",
	}

	ignorePatterns := cmd.Flags().StringSlice("ignore", nil, "Ingnore patterns as regular expressions")
	output := cmd.Flags().StringP("output", "o", "", "The file output path (if unset, will print to stdout)")
	packageName := cmd.Flags().StringP("package", "p", "static", "The package name for the generated file.")

	cmd.Run = func(_ *cobra.Command, args []string) {
		bundle := new(bindata.Bundle)
		bundle.PackageName = *packageName
		for _, ignorePattern := range *ignorePatterns {
			ignore, err := regexp.Compile(ignorePattern)
			if err != nil {
				fatal(err)
			}
			bundle.Ignores = append(bundle.Ignores, ignore)
		}

		var dst io.Writer
		if *output != "" {
			f, err := os.Create(*output)
			if err != nil {
				fatal(ex.New(err))
			}
			defer f.Close()
			dst = f
		} else {
			dst = os.Stdout
		}

		if err := bundle.Start(dst); err != nil {
			fatal(err)
		}

		for _, path := range args {
			if err := bundle.ProcessFile(dst, parsePathConfig(path)); err != nil {
				fatal(err)
			}
		}
		if err := bundle.Finish(dst); err != nil {
			fatal(err)
		}
	}
	fatal(cmd.Execute())
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n\n", err)
		os.Exit(1)
	}
}

// parseRecursive determines whether the given path has a recrusive indicator and
// returns a new path with the recursive indicator chopped off if it does.
//
//  ex:
//      /path/to/foo/...    -> (/path/to/foo, true)
//      /path/to/bar        -> (/path/to/bar, false)
func parsePathConfig(path string) bindata.PathConfig {
	if strings.HasSuffix(path, "/...") {
		return bindata.PathConfig{
			Path:      filepath.Clean(path[:len(path)-4]),
			Recursive: true,
		}
	}
	return bindata.PathConfig{
		Path:      filepath.Clean(path),
		Recursive: false,
	}
}
