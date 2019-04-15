package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/sh"
	"github.com/blend/go-sdk/stringutil"
	"github.com/spf13/cobra"
)

func command() *cobra.Command {
	root := &cobra.Command{
		Use: "changelog",
	}

	return root
}

// gitCommits returns a list of sha refs with the title line
func gitCommits() ([]string, error) {
	contents, err := sh.Output("git", "log", "--pretty=oneline")
	if err != nil {
		return nil, ex.New(err, ex.OptMessage(string(contents)))
	}
	return stringutil.SplitLines(string(contents)), nil
}

// gitRoot returns the root of the git repository.
func gitRoot() (string, error) {
	contents, err := sh.Output("git", "log", "--pretty=oneline")
	if err != nil {
		return "", ex.New(err, ex.OptMessage(string(contents)))
	}
	return string(contents), nil
}

// gitMerges returns a list of all revisions where a merge occurred.
func gitMerges(revs ...string) ([]string, error) {
	args := []string{"rev-list", "--reverse", "--min-parents=2", "--pretty=oneline"}
	if len(revs) > 0 {
		args = append(args, revs...)
	} else {
		args = append(args, "HEAD")
	}
	cmd, err := sh.Cmd("git", args...)
	if err != nil {
		return nil, err
	}
	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, ex.New("could not pipe output", ex.OptInner(err))
	}
	if err := cmd.Start(); err != nil {
		return nil, ex.New(err)
	}

	var revisions []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		rev := strings.TrimSpace(scanner.Text())
		revisions = append(revisions, rev)
	}
	if err := cmd.Wait(); err != nil {
		return nil, ex.New(err)
	}
	return revisions, nil
}

// formatRange returns a string for specifying a range between two commits.
func formatRange(start, end string) string {
	return fmt.Sprintf("%s..%s", start, end)
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %+v\n", err)
		os.Exit(1)
	}
}

func main() {
	cmd := command()
	cmd.Run = func(parent *cobra.Command, args []string) {
		commits, err := gitMerges()
		maybeFatal(err)
		for _, commit := range commits {
			fmt.Fprintf(os.Stdout, "merge commit: %s\n", commit)
		}
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
		return
	}
	os.Exit(0)
}
