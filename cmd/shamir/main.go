/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/shamir"
	"github.com/blend/go-sdk/stringutil"
)

func main() {
	root := &cobra.Command{
		Use:	"shamir",
		Short:	"shamir splits and combines secrets into a configurable number of parts",
	}
	root.AddCommand(NewSplitCommand(), NewCombineCommand())
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// NewSplitCommand returns a new split command.
func NewSplitCommand() *cobra.Command {
	var parts *int
	var threshold *int
	var inputFile *string
	var secret *string
	split := &cobra.Command{
		Use:	"split",
		Short:	"split takes a given input from stdin or a file and separates it into a configurable number of sections",
		Run: func(cmd *cobra.Command, args []string) {
			var contents []byte
			var err error
			if len(*secret) > 0 {
				contents = []byte(strings.TrimSpace(*secret))
			} else if strings.TrimSpace(*inputFile) == "-" {
				contents, err = io.ReadAll(os.Stdin)
			} else {
				contents, err = os.ReadFile(strings.TrimSpace(*inputFile))
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if len(contents) == 0 {
				fmt.Fprintln(os.Stderr, fmt.Errorf("invalid input; is empty"))
				os.Exit(1)
			}

			shares, err := shamir.Split(contents, *parts, *threshold)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			for _, share := range shares {
				fmt.Println(hex.EncodeToString(share))
			}
		},
	}
	secret = split.Flags().StringP("secret", "s", "", "the input secret")
	inputFile = split.Flags().StringP("file", "f", "-", "the input file ('-' instructs to read from stdin)")
	parts = split.Flags().IntP("parts", "p", 5, "the number of parts to split the secret into")
	threshold = split.Flags().IntP("threshold", "t", 2, "the number of parts required to form the original secret")
	return split
}

// NewCombineCommand returns a new combine command.
func NewCombineCommand() *cobra.Command {
	var input *string
	var parts *[]string
	combine := &cobra.Command{
		Use:	"combine",
		Short:	"combine takes a shard share and combines it into the final output",
		Run: func(cmd *cobra.Command, args []string) {
			var inputParts [][]byte
			var inputPartsEncoded []string
			if len(*parts) > 0 {
				inputPartsEncoded = *parts
			} else {
				var contents []byte
				var err error
				if strings.TrimSpace(*input) == "-" {
					contents, err = io.ReadAll(os.Stdin)
				} else {
					contents, err = os.ReadFile(strings.TrimSpace(*input))
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				inputPartsEncoded = stringutil.SplitLines(string(contents))
			}

			for _, part := range inputPartsEncoded {
				decoded, err := hex.DecodeString(strings.TrimSpace(part))
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				inputParts = append(inputParts, decoded)
			}

			original, err := shamir.Combine(inputParts)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println(string(original))
		},
	}
	input = combine.Flags().StringP("input", "i", "-", "the input file ('-' instructs to read from stdin)")
	parts = combine.Flags().StringArrayP("part", "p", nil, "individual parts to combine (must include the threshold amount")
	return combine
}
