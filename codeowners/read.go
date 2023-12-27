/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package codeowners

import (
	"bufio"
	"io"
	"strings"
)

// Read reads a codeowners file.
func Read(r io.Reader) (output File, err error) {
	scanner := bufio.NewScanner(r)
	var line string
	var codeownersEntry Source
	for scanner.Scan() {
		line = scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		} else if strings.HasPrefix(line, OwnersFileSourceComment) {
			codeownersEntry.Source = strings.TrimPrefix(line, OwnersFileSourceComment)
			continue
		} else if strings.HasPrefix(line, OwnersFileSourceEndComment) {
			output = append(output, codeownersEntry)
			codeownersEntry = Source{}
			continue
		} else if strings.HasPrefix(line, "#") {
			continue
		}

		var codeownersEntryPath Path
		codeownersEntryPath, err = ParsePath(line)
		if err != nil {
			return
		}
		codeownersEntry.Paths = append(codeownersEntry.Paths, codeownersEntryPath)
	}
	return
}
