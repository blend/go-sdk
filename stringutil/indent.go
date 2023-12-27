/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

import "strings"

// Indent applies an indent prefix to a given corpus.
func Indent(indent string, corpus string) string {
	return strings.Join(IndentLines(indent, SplitLines(corpus)), "\n")
}

// IndentLines adds a prefix to a given list of strings.
func IndentLines(indent string, corpus []string) []string {
	for index := 0; index < len(corpus); index++ {
		corpus[index] = indent + corpus[index]
	}
	return corpus
}
