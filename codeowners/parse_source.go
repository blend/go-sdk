/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package codeowners

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParseSource parses a source from a given path.
func ParseSource(repoRoot, sourcePath string) (*Source, error) {
	f, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	repoSourcePath, err := MakeRepositoryAbsolute(repoRoot, sourcePath)
	if err != nil {
		return nil, err
	}
	output := Source{
		Source: repoSourcePath,
	}

	var line string
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(strings.TrimSpace(line), "#") {	// ignore comments
			continue
		}
		if strings.Contains(line, "#") {
			return nil, fmt.Errorf("invalid codeowners file; must not contain inline comments")
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		pieces := strings.Fields(line)
		if len(pieces) < 2 {
			return nil, fmt.Errorf("invalid codeowners file; must contain a path glob and an owner on a single line")
		}
		pathGlob, err := MakeRepositoryAbsolute(repoRoot, filepath.Join(filepath.Dir(sourcePath), pieces[0]))
		if err != nil {
			return nil, err
		}
		output.Paths = append(output.Paths, Path{
			PathGlob:	pathGlob,
			Owners:		pieces[1:],
		})
	}
	return &output, nil
}
