/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"fmt"
	"strings"
)

// Append appends errors together, creating a multi-error.
func Append(err error, errs ...error) error {
	if len(errs) == 0 {
		return err
	}
	var all []error
	if err != nil {
		if me, ok := err.(Multi); ok {
			all = me
		} else {
			all = append(all, NewWithStackDepth(err, DefaultNewStartDepth+1))
		}
	}
	for _, extra := range errs {
		if extra != nil {
			if _, ok := extra.(Multi); !ok {
				extra = NewWithStackDepth(extra, DefaultNewStartDepth+1)
			}
			all = append(all, extra)
		}
	}
	if len(all) == 0 {
		return nil
	}
	if len(all) == 1 {
		return all[0]
	}
	return Multi(all)
}

// Multi represents an array of errors.
type Multi []error

// Unwrap returns all the errors in the multi error (basically itself)
func (m Multi) Unwrap() []error {
	return m
}

// Error implements error.
func (m Multi) Error() string {
	formatted, _ := m.errorString(10, 5, 0)
	return formatted
}

// FullError returns the full error message.
func (m Multi) FullError() string {
	formatted, _ := m.errorString(-1, -1, 0)
	return formatted
}

// errorString returns the error string with a length limit and a depth limit,
// along with the total number of errors in the Multi error tree.
// -1 means no limit.
func (m Multi) errorString(listLengthLimit, depthLimit, depth int) (string, int) {
	if len(m) == 0 {
		return "", 0
	}
	prefix := "\t"
	if depthLimit >= 0 && depth+1 > depthLimit {
		total := countErrors(m)
		return fmt.Sprintf("%s\n%s... depth limit reached ...", m.header(total), prefix), total
	}

	total := 0
	var points []string
	for i, err := range m {
		if i == listLengthLimit {
			expanded := total
			for _, err := range m[i:] {
				total += countErrors(err)
			}
			points = append(points, fmt.Sprintf("... and %d more", total-expanded))
			break
		}
		if me, ok := err.(Multi); ok {
			formatted, count := me.errorString(listLengthLimit, depthLimit, depth+1)
			points = append(points, fmt.Sprintf("* %s", indent(prefix+prefix, formatted)))
			total += count
			continue
		}
		points = append(points, fmt.Sprintf("* %s", indent(prefix+prefix, err.Error())))
		total++
	}

	return fmt.Sprintf("%s\n%s%s", m.header(total), prefix, strings.Join(points, "\n"+prefix)), total
}

func (Multi) header(total int) string {
	if total == 1 {
		return "1 error occurred:"
	}
	return fmt.Sprintf("%d errors occurred:", total)
}

func countErrors(err error) int {
	if me, ok := err.(Multi); ok {
		count := 0
		for _, err := range me {
			count += countErrors(err)
		}
		return count
	}
	return 1
}

// Indent functions from https://git.blendlabs.com/blend/go/blob/ad2d96d6f0d62d9a2a2037945663f44d438f6300/sdk/stringutil/indent.go
// to avoid import cycle.

// indent applies an indent prefix to a given corpus except the first line.
func indent(prefix, corpus string) string {
	return strings.Join(indentLines(prefix, strings.Split(corpus, "\n")), "\n")
}

// indentLines adds a prefix to a given list of strings except the first string.
func indentLines(prefix string, corpus []string) []string {
	for index := 1; index < len(corpus); index++ {
		corpus[index] = prefix + corpus[index]
	}
	return corpus
}
