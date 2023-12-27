/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

// SplitLinesOptions are options for the SplitLines function.
type SplitLinesOptions struct {
	IncludeNewline		bool
	IncludeEmptyLines	bool
}

// SplitLinesOption is a mutator for SplitLinesOptions.
type SplitLinesOption func(*SplitLinesOptions)

// OptSplitLinesIncludeNewLine sets if we should omit newlines in the returned lines.
func OptSplitLinesIncludeNewLine(include bool) SplitLinesOption {
	return func(opts *SplitLinesOptions) {
		opts.IncludeNewline = include
	}
}

// OptSplitLinesIncludeEmptyLines sets if we should omit newlines in the returned lines.
func OptSplitLinesIncludeEmptyLines(include bool) SplitLinesOption {
	return func(opts *SplitLinesOptions) {
		opts.IncludeEmptyLines = include
	}
}

// SplitLines splits a corpus into individual lines by the ascii control character `\n`.
// You can control some behaviors of the splitting process with variadic options.
func SplitLines(contents string, opts ...SplitLinesOption) []string {
	contentRunes := []rune(contents)

	var options SplitLinesOptions
	for _, opt := range opts {
		opt(&options)
	}

	var output []string

	const newline = '\n'

	var line []rune
	var c rune
	for index := 0; index < len(contentRunes); index++ {
		c = contentRunes[index]

		// if we hit a newline
		if c == newline {

			// if we should omit newlines
			if options.IncludeNewline {
				line = append(line, c)
			}

			// if we should omit empty lines
			if options.IncludeNewline {
				if len(line) == 1 && !options.IncludeEmptyLines {
					line = nil
					continue
				}
			} else {
				if len(line) == 0 && !options.IncludeEmptyLines {
					line = nil
					continue
				}
			}

			// add to the output
			output = append(output, string(line))
			line = nil
			continue
		}

		// add non-newline characters to the line
		line = append(line, c)
		continue
	}

	// add anything left
	if len(line) > 0 {
		output = append(output, string(line))
	}
	return output
}
