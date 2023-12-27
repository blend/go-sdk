/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package codeowners

import (
	"bufio"
	"bytes"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// ParseGoComments parses a files comments.
func ParseGoComments(repoRoot, sourcePath, linePrefix string) (*Source, error) {
	contents, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, sourcePath, contents, parser.ImportsOnly|parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var owners []string
	var corpus, line string
	for _, commentGroup := range fileAst.Comments {
		corpus = commentGroup.Text()
		scanner := bufio.NewScanner(bytes.NewBufferString(corpus))
		// scan the corpus lines for `//github:codeowners prefixes`
		for scanner.Scan() {
			line = strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, linePrefix) {
				owners = append(owners, strings.Fields(strings.TrimPrefix(line, linePrefix))...)
			}
		}
	}

	if len(owners) == 0 {
		return nil, nil
	}
	repoSourcePath, err := MakeRepositoryAbsolute(repoRoot, sourcePath)
	if err != nil {
		return nil, err
	}
	pathGlob, err := MakeRepositoryAbsolute(repoRoot, sourcePath)
	if err != nil {
		return nil, err
	}
	return &Source{
		Source:	repoSourcePath,
		Paths: []Path{
			{
				PathGlob:	pathGlob,
				Owners:		owners,
			},
		},
	}, nil
}
