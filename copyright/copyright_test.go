/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package copyright

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ref"
)

func Test_Copyright_GetStdout(t *testing.T) {
	its := assert.New(t)

	c := New()

	its.Equal(os.Stdout, c.GetStdout())
	buf := new(bytes.Buffer)
	c.Stdout = buf
	its.Equal(c.Stdout, c.GetStdout())
	c.Quiet = ref.Bool(true)
	its.Equal(io.Discard, c.GetStdout())
}

func Test_Copyright_GetStderr(t *testing.T) {
	its := assert.New(t)

	c := New()

	its.Equal(os.Stderr, c.GetStderr())
	buf := new(bytes.Buffer)
	c.Stderr = buf
	its.Equal(buf, c.GetStderr())
	c.Quiet = ref.Bool(true)
	its.Equal(io.Discard, c.GetStderr())
}

func Test_Copyright_mergeFileSections(t *testing.T) {
	its := assert.New(t)

	merged := Copyright{}.mergeFileSections([]byte("foo"), []byte("bar"), []byte("baz"))
	its.Equal("foobarbaz", string(merged))
}

func Test_Copyright_fileHasCopyrightHeader(t *testing.T) {
	its := assert.New(t)

	var goodCorpus = []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(OptYear(2021))
	its.Nil(err)

	goodCorpusWithNotice := Copyright{}.mergeFileSections([]byte(notice), goodCorpus)
	its.Contains(string(goodCorpusWithNotice), "Copyright (c) 2021")
	its.True((Copyright{}).fileHasCopyrightHeader(goodCorpusWithNotice, []byte(notice)))
}

func Test_Copyright_fileHasCopyrightHeader_invalid(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	var invalidCorpus = []byte(`foo
bar
baz
`)
	expectedNotice, err := generateGoNotice(OptYear(2021))
	its.Nil(err)

	its.False(c.fileHasCopyrightHeader(invalidCorpus, []byte(expectedNotice)), "we haven't added the notice")
}

func Test_Copyright_fileHasCopyrightHeader_differentYear(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	var goodCorpus = []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(OptYear(2020))
	its.Nil(err)

	goodCorpusWithNotice := c.mergeFileSections(notice, goodCorpus)
	its.Contains(string(goodCorpusWithNotice), "Copyright (c) 2020")

	newNotice, err := generateGoNotice(OptYear(2021))
	its.Nil(err)

	its.True(c.fileHasCopyrightHeader(goodCorpusWithNotice, []byte(newNotice)))
}

func Test_Copyright_fileHasCopyrightHeader_leadingWhitespace(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	var goodCorpus = []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(OptYear(2021))
	its.Nil(err)

	goodCorpusWithNotice := c.mergeFileSections([]byte("\n\n"), notice, goodCorpus)
	its.HasPrefix(string(goodCorpusWithNotice), "\n\n")
	its.Contains(string(goodCorpusWithNotice), "Copyright (c) 2021")

	its.True(c.fileHasCopyrightHeader(goodCorpusWithNotice, []byte(notice)))
}

func Test_Copyright_goBuildTagMatch(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	buildTag := []byte(`// +build foo

`)
	corpus := []byte(`foo
bar
baz
`)

	file := (Copyright{}).mergeFileSections(buildTag, corpus)

	its.False(goBuildTagMatch.Match(corpus))
	its.True(goBuildTagMatch.Match(c.mergeFileSections(buildTag)))

	found := goBuildTagMatch.FindAll(file, -1)
	its.NotEmpty(found)
	its.True(goBuildTagMatch.Match(file))
}

func Test_Copyright_goBuildTagsMatch(t *testing.T) {
	its := assert.New(t)

	file := []byte(goBuildTags1)	// testutil.GetTestFixture(its, "buildtags1.go")
	its.True(goBuildTagMatch.Match(file))
	found := goBuildTagMatch.Find(file)
	its.Equal("//go:build tag1\n// +build tag1\n\n", string(found))

	file2 := []byte(goBuildTags2)	// testutil.GetTestFixture(its, "buildtags2.go")
	its.True(goBuildTagMatch.Match(file2))
	found2 := goBuildTagMatch.Find(file2)

	expected := `// +build tag5
//go:build tag1 && tag2 && tag3
// +build tag1,tag2,tag3
// +build tag6

`
	its.Equal(expected, string(found2))

	file3 := []byte(goBuildTags3)	// testutil.GetTestFixture(its, "buildtags3.go")
	its.True(goBuildTagMatch.Match(file3))
	found3 := goBuildTagMatch.Find(file3)
	its.Equal("//go:build tag1 & tag2\n\n", string(found3))
}

func Test_Copyright_goInjectNotice(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	file := []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(OptYear(2021))
	its.Nil(err)

	output := c.goInjectNotice("foo.go", file, notice)
	its.Contains(string(output), "Copyright (c) 2021")
	its.HasSuffix(string(output), string(file))
}

func Test_Copyright_goInjectNotice_buildTag(t *testing.T) {
	its := assert.New(t)
	c := Copyright{}

	buildTag := []byte(`// +build foo`)
	corpus := []byte(`foo
bar
baz
`)

	file := c.mergeFileSections(buildTag, []byte("\n\n"), corpus)

	notice, err := generateGoNotice(OptYear(2021))
	its.Nil(err)

	output := c.goInjectNotice("foo.go", file, notice)
	its.Contains(string(output), "Copyright (c) 2021")
	its.HasPrefix(string(output), string(buildTag)+"\n")
	its.HasSuffix(string(output), string(corpus))

	outputRepeat := c.goInjectNotice("foo.go", output, notice)
	its.Empty(outputRepeat, "inject notice functions should return an empty slice if the header already exists")
}

func Test_Copyright_goInjectNotice_goBuildTags(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Name	string
		Input	string
		Expect	string
	}

	cases := []testCase{
		{
			Name:	"standard build tags",
			Input:	goBuildTags1,	// "buildtags1.go",
			Expect:	goldenGoBuildTags1,
		},
		{
			Name:	"multiple build tags",
			Input:	goBuildTags2,	// "buildtags2.go",
			Expect:	goldenGoBuildTags2,
		},
		{
			Name:	"build tags split across file",
			Input:	goBuildTags3,	// "buildtags3.go",
			Expect:	goldenGoBuildTags3,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			it := assert.New(t)
			c := Copyright{}

			notice, err := generateGoNotice(OptYear(2001))
			it.Nil(err)

			output := c.goInjectNotice("foo.go", []byte(tc.Input), notice)
			it.Equal(string(output), tc.Expect)	// testutil.AssertGoldenFile(it, output, tc.TestFile)

			outputRepeat := c.goInjectNotice("foo.go", output, notice)
			it.Empty(outputRepeat)
		})
	}
}

func Test_Copyright_tsInjectNotice_tsReferenceTags(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Name	string
		Input	string
		Expect	string
	}

	cases := []testCase{
		{
			Name:	"single reference tag",
			Input:	tsReferenceTag,
			Expect:	goldenTsReferenceTag,
		},
		{
			Name:	"multiple reference tags",
			Input:	tsReferenceTags,
			Expect:	goldenTsReferenceTags,
		},
		{
			Name:	"no reference tags",
			Input:	tsTest,	// "buildtags3.go",
			Expect:	goldenTs,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			it := assert.New(t)
			c := Copyright{}

			notice, err := generateTypescriptNotice(OptYear(2022))
			it.Nil(err)

			output := c.tsInjectNotice("foo.ts", []byte(tc.Input), notice)
			it.Equal(tc.Expect, string(output))

			outputRepeat := c.tsInjectNotice("foo.ts", output, notice)
			it.Empty(outputRepeat)
		})
	}
}

func Test_Copyright_injectNotice_typescript(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	file := []byte(`foo
bar
baz
`)

	notice, err := generateTypescriptNotice(OptYear(2001))
	its.Nil(err)

	output := c.injectNotice("foo.ts", file, notice)
	its.Contains(string(output), "Copyright (c) 2001")
	its.HasSuffix(string(output), string(file))

	outputRepeat := c.injectNotice("foo.ts", output, notice)
	its.Empty(outputRepeat, "inject notice functions should return an empty slice if the header already exists")
}

func Test_Copyright_injectNotice_typescript_referenceTags(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	file := []byte(tsReferenceTags)

	notice, err := generateTypescriptNotice(OptYear(2001))
	its.Nil(err)

	output := c.injectNotice("foo.ts", file, notice)
	its.Contains(string(output), "Copyright (c) 2001")
	its.HasSuffix(string(output), string(file))

	outputRepeat := c.injectNotice("foo.ts", output, notice)
	its.Empty(outputRepeat, "inject notice functions should return an empty slice if the header already exists")
}

func Test_Copyright_goInjectNotice_openSource(t *testing.T) {
	its := assert.New(t)

	c := new(Copyright)

	file := []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(
		OptYear(2021),
		OptLicense("Apache 2.0"),
		OptRestrictions(DefaultRestrictionsOpenSource),
	)
	its.Nil(err)

	output := c.goInjectNotice("foo.go", file, notice)
	its.Contains(string(output), "Copyright (c) 2021")
	its.Contains(string(output), "Use of this source code is governed by a Apache 2.0")
	its.HasSuffix(string(output), string(file))
}

func generateGoNotice(opts ...Option) ([]byte, error) {
	c := New(opts...)
	noticeBody, err := c.compileNoticeBodyTemplate(c.NoticeBodyTemplateOrDefault())
	if err != nil {
		return nil, err
	}

	compiled, err := c.compileNoticeTemplate(goNoticeTemplate, noticeBody)
	if err != nil {
		return nil, err
	}
	return []byte(compiled), nil
}

func generateTypescriptNotice(opts ...Option) ([]byte, error) {
	c := New(opts...)
	noticeBody, err := c.compileNoticeBodyTemplate(c.NoticeBodyTemplateOrDefault())
	if err != nil {
		return nil, err
	}

	compiled, err := c.compileNoticeTemplate(tsNoticeTemplate, noticeBody)
	if err != nil {
		return nil, err
	}
	return []byte(compiled), nil
}

func Test_Copyright_GetNoticeTemplate(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	noticeTemplate, ok := c.noticeTemplateByExtension(".js")
	its.True(ok)
	its.Equal(jsNoticeTemplate, noticeTemplate)

	// it handles no dot prefix
	noticeTemplate, ok = c.noticeTemplateByExtension("js")
	its.True(ok)
	its.Equal(jsNoticeTemplate, noticeTemplate)

	// it handles another file type
	noticeTemplate, ok = c.noticeTemplateByExtension(".go")
	its.True(ok)
	its.Equal(goNoticeTemplate, noticeTemplate)

	noticeTemplate, ok = c.noticeTemplateByExtension("not-a-real-extension")
	its.False(ok)
	its.Empty(noticeTemplate)

	withDefault := Copyright{
		Config: Config{
			FallbackNoticeTemplate: "this is just a test",
		},
	}

	noticeTemplate, ok = withDefault.noticeTemplateByExtension("not-a-real-extension")
	its.True(ok)
	its.Equal("this is just a test", noticeTemplate)
}

type mockInfoDir string

func (mid mockInfoDir) Name() string		{ return string(mid) }
func (mid mockInfoDir) Size() int64		{ return 1 << 8 }
func (mid mockInfoDir) Mode() fs.FileMode	{ return fs.FileMode(0755) }
func (mid mockInfoDir) ModTime() time.Time	{ return time.Now().UTC() }
func (mid mockInfoDir) IsDir() bool		{ return true }
func (mid mockInfoDir) Sys() interface{}	{ return nil }

type mockInfoFile string

func (mif mockInfoFile) Name() string		{ return string(mif) }
func (mif mockInfoFile) Size() int64		{ return 1 << 8 }
func (mif mockInfoFile) Mode() fs.FileMode	{ return fs.FileMode(0755) }
func (mif mockInfoFile) ModTime() time.Time	{ return time.Now().UTC() }
func (mif mockInfoFile) IsDir() bool		{ return false }
func (mif mockInfoFile) Sys() interface{}	{ return nil }

func Test_Copyright_includeOrExclude(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	testCases := [...]struct {
		Config		Config
		Path		string
		Info		fs.FileInfo
		Expected	error
	}{
		/*0*/ {Config: Config{}, Path: ".", Info: mockInfoDir("."), Expected: ErrWalkSkip},
		/*1*/ {Config: Config{Excludes: []string{"/foo/**"}}, Path: "/foo/bar", Info: mockInfoDir("bar"), Expected: filepath.SkipDir},
		/*2*/ {Config: Config{Excludes: []string{"/foo/**"}}, Path: "/foo/bar/baz.jpg", Info: mockInfoFile("baz.jpg"), Expected: ErrWalkSkip},
		/*3*/ {Config: Config{IncludeFiles: []string{"/foo/bar/*.jpg"}}, Path: "/foo/bar/baz.jpg", Info: mockInfoFile("baz.jpg"), Expected: nil},
		/*4*/ {Config: Config{Excludes: []string{}, IncludeFiles: []string{}}, Path: "/foo/bar/baz.jpg", Info: mockInfoFile("baz.jpg"), Expected: ErrWalkSkip},
		/*5*/ {Config: Config{}, Path: "/foo/bar/baz.jpg", Info: mockInfoFile("baz.jpg"), Expected: nil},
		/*6*/ {Config: Config{}, Path: "/foo/bar/baz.jpg", Info: mockInfoDir("baz"), Expected: ErrWalkSkip},
	}

	for index, tc := range testCases {
		c := Copyright{Config: tc.Config}
		its.Equal(tc.Expected, c.includeOrExclude(".", tc.Path, tc.Info), fmt.Sprintf("test %d", index))
	}
}

const (
	tsFile0	= `import * as axios from 'axios';`
	tsFile1	= `/// <reference path="../types/testing.d.ts" />
/// <reference path="../types/something.d.ts" />
/// <reference path="../types/somethingElse.d.ts" />
/// <reference path="../types/somethingMore.d.ts" />
/// <reference path="../types/somethingLess.d.ts" />

	import * as axios from 'axios';
`

	pyFile0	= `from __future__ import print_function
	
		import logging
		import os
		import shutil
		import sys
		import requests
		import uuid
		import json`

	goFile0	= `// +build tools
		package tools
		
		import (
			// goimports organizes imports for us
			_ "golang.org/x/tools/cmd/goimports"
		
			// revive is an opinionated linter
			_ "github.com/mgechev/revive"
		
			// ineffassign is an opinionated linter
			_ "github.com/gordonklaus/ineffassign"
		
			// staticcheck is ineffassign but better
			_ "honnef.co/go/tools/cmd/staticcheck"
		)
		`
)

// createTestFS creates a temp dir with files in them, with _no_ copyright headers.
//
// there should be at least (1) failure.
func createTestFS(its *assert.Assertions) (tempDir string, revert func()) {
	// create a temp dir
	var err error
	tempDir, err = os.MkdirTemp("", "copyright_test")
	its.Nil(err)
	revert = func() {
		os.RemoveAll(tempDir)
	}

	// create some files
	err = os.MkdirAll(filepath.Join(tempDir, "foo", "bar"), 0755)
	its.Nil(err)
	err = os.MkdirAll(filepath.Join(tempDir, "bar", "foo"), 0755)
	its.Nil(err)

	err = os.MkdirAll(filepath.Join(tempDir, "not-bar", "not-foo"), 0755)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "file0.py"), []byte(pyFile0), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "file0.ts"), []byte(tsFile0), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "file1.ts"), []byte(tsFile1), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "file0.go"), []byte(goFile0), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "foo", "bar", "file0.py"), []byte(pyFile0), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "foo", "bar", "file0.ts"), []byte(tsFile0), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "foo", "bar", "file1.ts"), []byte(tsFile1), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "foo", "bar", "file0.go"), []byte(goFile0), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "bar", "foo", "file0.py"), []byte(pyFile0), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "bar", "foo", "file0.ts"), []byte(tsFile0), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "bar", "foo", "file1.ts"), []byte(tsFile1), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "bar", "foo", "file0.go"), []byte(goFile0), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "not-bar", "not-foo", "file0.py"), []byte(pyFile0), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "not-bar", "not-foo", "file0.ts"), []byte(tsFile0), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "not-bar", "not-foo", "file1.ts"), []byte(tsFile1), 0644)
	its.Nil(err)

	err = os.WriteFile(filepath.Join(tempDir, "not-bar", "not-foo", "file0.go"), []byte(goFile0), 0644)
	its.Nil(err)
	return
}

func Test_Copyright_Walk(t *testing.T) {
	its := assert.New(t)

	tempDir, revert := createTestFS(its)
	defer revert()

	c := New(
		OptIncludeFiles("*.py", "*.ts"),
		OptExcludes("*/not-bar/*", "*/not-foo/*"),
	)

	var err error
	var seen []string
	err = c.Walk(context.TODO(), func(path string, info os.FileInfo, file, notice []byte) error {
		seen = append(seen, path)
		return nil
	}, tempDir)
	its.Nil(err)
	expected := []string{
		filepath.Join(tempDir, "bar", "foo", "file0.py"),
		filepath.Join(tempDir, "bar", "foo", "file0.ts"),
		filepath.Join(tempDir, "bar", "foo", "file1.ts"),
		filepath.Join(tempDir, "file0.py"),
		filepath.Join(tempDir, "file0.ts"),
		filepath.Join(tempDir, "file1.ts"),
		filepath.Join(tempDir, "foo", "bar", "file0.py"),
		filepath.Join(tempDir, "foo", "bar", "file0.ts"),
		filepath.Join(tempDir, "foo", "bar", "file1.ts"),
	}
	its.Equal(expected, seen)
}

func Test_Copyright_Walk_noExitFirst(t *testing.T) {
	its := assert.New(t)

	tempDir, revert := createTestFS(its)
	defer revert()

	c := New(
		OptIncludeFiles("*.py", "*.ts"),
		OptExcludes("*/not-bar/*", "*/not-foo/*"),
		OptExitFirst(false),
	)

	var err error
	var seen []string
	err = c.Walk(context.TODO(), func(path string, info os.FileInfo, file, notice []byte) error {
		seen = append(seen, path)
		if len(seen) > 0 {
			return ErrFailure
		}
		return nil
	}, tempDir)
	its.NotNil(err)
	expected := []string{
		filepath.Join(tempDir, "bar", "foo", "file0.py"),
		filepath.Join(tempDir, "bar", "foo", "file0.ts"),
		filepath.Join(tempDir, "bar", "foo", "file1.ts"),
		filepath.Join(tempDir, "file0.py"),
		filepath.Join(tempDir, "file0.ts"),
		filepath.Join(tempDir, "file1.ts"),
		filepath.Join(tempDir, "foo", "bar", "file0.py"),
		filepath.Join(tempDir, "foo", "bar", "file0.ts"),
		filepath.Join(tempDir, "foo", "bar", "file1.ts"),
	}
	its.Equal(expected, seen)
}

func Test_Copyright_Walk_exitFirst(t *testing.T) {
	its := assert.New(t)

	tempDir, revert := createTestFS(its)
	defer revert()

	c := New(
		OptIncludeFiles("*.py", "*.ts"),
		OptExcludes("*/not-bar/*", "*/not-foo/*"),
		OptExitFirst(true),
	)

	var err error
	var seen []string
	err = c.Walk(context.TODO(), func(path string, info os.FileInfo, file, notice []byte) error {
		seen = append(seen, path)
		if len(seen) > 0 {
			return ErrFailure
		}
		return nil
	}, tempDir)
	its.NotNil(err)
	expected := []string{
		filepath.Join(tempDir, "bar", "foo", "file0.py"),
	}
	its.Equal(expected, seen)
}

func Test_Copyright_Inject(t *testing.T) {
	its := assert.New(t)

	tempDir, revert := createTestFS(its)
	defer revert()

	c := New(
		OptIncludeFiles("*.py", "*.ts", "*.go"),
		OptExcludes("*/not-bar/*", "*/not-foo/*"),
	)

	err := c.Inject(context.TODO(), tempDir)
	its.Nil(err)

	contents, err := os.ReadFile(filepath.Join(tempDir, "bar", "foo", "file0.py"))
	its.Nil(err)
	its.HasPrefix(string(contents), "#\n# Copyright")

	contents, err = os.ReadFile(filepath.Join(tempDir, "bar", "foo", "file0.ts"))
	its.Nil(err)
	its.HasPrefix(string(contents), "/**\n * Copyright")
}

func Test_Copyright_Inject_Shebang(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	tempDir, err := os.MkdirTemp("", "copyright_test")
	its.Nil(err)
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	// Write `shift.py` without
	contents := strings.Join([]string{
		"\r\t",
		"  #!/usr/bin/env python",
		"",
		"def main():",
		`    print("Hello world")`,
		"",
		"",
		`if __name__ == "__main__":`,
		"    main()",
		"",
	}, "\n")
	filename := filepath.Join(tempDir, "shift.py")
	err = os.WriteFile(filename, []byte(contents), 0755)
	its.Nil(err)

	// Actually inject
	c := New(OptIncludeFiles("*shift.py"))
	err = c.Inject(context.TODO(), tempDir)
	its.Nil(err)

	// Verify injected contents are as expected
	contentInjected, err := os.ReadFile(filename)
	its.Nil(err)
	expected := strings.Join([]string{
		"\r\t",
		"  #!/usr/bin/env python",
		"#",
		"# " + expectedNoticePrefix(its),
		"# " + DefaultRestrictionsInternal,
		"#",
		"",
		"",
		"def main():",
		`    print("Hello world")`,
		"",
		"",
		`if __name__ == "__main__":`,
		"    main()",
		"",
	}, "\n")
	its.Equal(expected, string(contentInjected))

	// Verify no-op if notice header is already present
	err = c.Inject(context.TODO(), tempDir)
	its.Nil(err)
	contentInjected, err = os.ReadFile(filename)
	its.Nil(err)
	its.Equal(expected, string(contentInjected))
}

func Test_Copyright_Verify(t *testing.T) {
	its := assert.New(t)

	tempDir, revert := createTestFS(its)
	defer revert()

	c := New(
		OptIncludeFiles("*.py", "*.ts", "*.go"),
		OptExcludes("*/not-bar/*", "*/not-foo/*"),
		OptExitFirst(false),
	)
	c.Stdout = new(bytes.Buffer)
	c.Stderr = new(bytes.Buffer)

	err := c.Verify(context.TODO(), tempDir)
	its.NotNil(err, "we must record a failure from walking the test fs")

	err = c.Inject(context.TODO(), tempDir)
	its.Nil(err)

	err = c.Verify(context.TODO(), tempDir)
	its.Nil(err)
}

func Test_Copyright_Verify_Shebang(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	tempDir, err := os.MkdirTemp("", "copyright_test")
	its.Nil(err)
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	// Write `shift.py` already injected
	contents := strings.Join([]string{
		"#!/usr/bin/env python",
		"#",
		"# " + expectedNoticePrefix(its),
		"# " + DefaultRestrictionsInternal,
		"#",
		"",
		"",
		"def main():",
		`    print("Hello world")`,
		"",
		"",
		`if __name__ == "__main__":`,
		"    main()",
		"",
	}, "\n")
	filename := filepath.Join(tempDir, "shift.py")
	err = os.WriteFile(filename, []byte(contents), 0755)
	its.Nil(err)

	// Verify present
	cfg := Config{
		ShowDiff:	ref.Bool(false),
		Quiet:		ref.Bool(true),
		IncludeFiles:	[]string{"*shift.py"},
	}
	c := New(OptConfig(cfg))
	err = c.Verify(context.TODO(), tempDir)
	its.Nil(err)

	// Write without and fail
	contents = strings.Join([]string{
		"#!/usr/bin/env python",
		"def main():",
		`    print("Hello world")`,
		"",
		"",
		`if __name__ == "__main__":`,
		"    main()",
		"",
	}, "\n")
	err = os.WriteFile(filename, []byte(contents), 0755)
	its.Nil(err)
	err = c.Verify(context.TODO(), tempDir)
	its.Equal("failure; one or more steps failed", fmt.Sprintf("%v", err))
}

func Test_Copyright_Remove(t *testing.T) {
	its := assert.New(t)

	tempDir, revert := createTestFS(its)
	defer revert()

	c := New(
		OptIncludeFiles("*.py", "*.ts", "*.go"),
		OptExcludes("*/not-bar/*", "*/not-foo/*"),
	)
	c.Stdout = new(bytes.Buffer)
	c.Stderr = new(bytes.Buffer)

	err := c.Inject(context.TODO(), tempDir)
	its.Nil(err)

	err = c.Verify(context.TODO(), tempDir)
	its.Nil(err)

	err = c.Remove(context.TODO(), tempDir)
	its.Nil(err)

	err = c.Verify(context.TODO(), tempDir)
	its.NotNil(err)
}

func Test_Copyright_Remove_Shebang(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	tempDir, err := os.MkdirTemp("", "copyright_test")
	its.Nil(err)
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	// Write `shift.py` already injected
	contents := strings.Join([]string{
		"#!/usr/bin/env python",
		"#",
		"# " + expectedNoticePrefix(its),
		"# " + DefaultRestrictionsInternal,
		"#",
		"",
		"",
		"def main():",
		`    print("Hello world")`,
		"",
		"",
		`if __name__ == "__main__":`,
		"    main()",
		"",
	}, "\n")
	filename := filepath.Join(tempDir, "shift.py")
	err = os.WriteFile(filename, []byte(contents), 0755)
	its.Nil(err)

	// Actually remove
	c := New(OptIncludeFiles("*shift.py"))
	err = c.Remove(context.TODO(), tempDir)
	its.Nil(err)

	// Verify removed contents are as expected
	contentRemoved, err := os.ReadFile(filename)
	its.Nil(err)
	expected := strings.Join([]string{
		"#!/usr/bin/env python",
		"def main():",
		`    print("Hello world")`,
		"",
		"",
		`if __name__ == "__main__":`,
		"    main()",
		"",
	}, "\n")
	its.Equal(expected, string(contentRemoved))

	// Verify no-op if notice header is already removed
	err = c.Remove(context.TODO(), tempDir)
	its.Nil(err)
	contentRemoved, err = os.ReadFile(filename)
	its.Nil(err)
	its.Equal(expected, string(contentRemoved))
}

func Test_Copyright_Walk_singleFileRoot(t *testing.T) {
	its := assert.New(t)

	tempDir, revert := createTestFS(its)
	defer revert()

	c := New(
		OptIncludeFiles("*.py", "*.ts"),
		OptExcludes("*/not-bar/*", "*/not-foo/*"),
	)

	var err error
	var seen []string
	err = c.Walk(context.TODO(), func(path string, info os.FileInfo, file, notice []byte) error {
		seen = append(seen, path)
		return nil
	}, filepath.Join(tempDir, "file0.py"))
	its.Nil(err)
	expected := []string{
		filepath.Join(tempDir, "file0.py"),
	}
	its.Equal(expected, seen)
}

func expectedNoticePrefix(its *assert.Assertions) string {
	vars := map[string]string{
		"Year":		fmt.Sprintf("%d", time.Now().UTC().Year()),
		"Company":	DefaultCompany,
		"Restrictions":	"",
	}
	tmpl := template.New("output")
	_, err := tmpl.Parse(DefaultNoticeBodyTemplate)
	its.Nil(err)
	prefixBuffer := new(bytes.Buffer)
	err = tmpl.Execute(prefixBuffer, vars)
	its.Nil(err)
	return strings.TrimRight(prefixBuffer.String(), "\n")
}
