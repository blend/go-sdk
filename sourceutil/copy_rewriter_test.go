/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sourceutil

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/blend/go-sdk/ref"
)

func Test_CopyRewrite_rewriteGoAst(t *testing.T) {
	t.Parallel()

	file := `package foo

import "fmt"
import "foo/bar"
import foo "bar/foo"
import "sdk/httpmetrics"

func main() {
	fmt.Println(bar.Foo)
	println(foo.Bar)
}
`

	expected := `package foo

import "fmt"
import "golang.org/foo/bar"
import foo "golang.org/bar/foo"
import httpmetrics "golang.org/sdk/stats/httpmetrics"

func main() {
	fmt.Printf(bar.Foo)
	println(foo.Bar)
}
`

	ctx := context.Background()
	cr := CopyRewriter{
		GoImportVisitors: []GoImportVisitor{
			GoImportRewritePrefix("foo", "golang.org/foo"),
			GoImportRewritePrefix("bar", "golang.org/bar"),
			GoImportRewrite(
				OptGoImportPathMatches("^sdk/httpmetrics"),
				OptGoImportAddName("httpmetrics"),
				OptGoImportSetPath("golang.org/sdk/stats/httpmetrics"),
			),
		},
		GoAstVistiors: []GoAstVisitor{
			GoAstRewrite(
				GoIsPackageCall("fmt", "Println"),
				GoRewritePackageCall("fmt", "Printf"),
			),
		},
		Debug:	ref.Bool(true),
	}
	var buf bytes.Buffer
	if err := cr.rewriteGoAst(ctx, "main.go", []byte(file), &buf); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if buf.String() == "" {
		t.Errorf("buffer was empty")
		t.FailNow()
	}

	if !strings.HasPrefix(buf.String(), expected) {
		t.Logf(buf.String())
		t.Errorf("invalid output")
		t.FailNow()
	}
}
