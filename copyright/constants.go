/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package copyright

import (
	"errors"
	"regexp"
)

// DefaultCompany is the default company to inject into the notice template.
const DefaultCompany = "Blend Labs, Inc"

// DefaultOpenSourceLicense is the default open source license.
const DefaultOpenSourceLicense = "MIT"

// DefaultRestrictionsInternal are the default copyright restrictions to inject into the notice template.
const DefaultRestrictionsInternal = "Blend Confidential - Restricted"

// DefaultRestrictionsOpenSource are the default open source restrictions.
const DefaultRestrictionsOpenSource = `Use of this source code is governed by a {{ .License }} license that can be found in the LICENSE file.`

// DefaultNoticeBodyTemplate is the default notice body template.
const DefaultNoticeBodyTemplate = `Copyright (c) {{ .Year }} - Present. {{ .Company }}. All rights reserved
{{ .Restrictions }}`

// Extension
const (
	ExtensionUnknown = ""
	ExtensionCSS     = ".css"
	ExtensionGo      = ".go"
	ExtensionHTML    = ".html"
	ExtensionJS      = ".js"
	ExtensionJSX     = ".jsx"
	ExtensionPy      = ".py"
	ExtensionSASS    = ".sass"
	ExtensionSCSS    = ".scss"
	ExtensionTS      = ".ts"
	ExtensionTSX     = ".tsx"
	ExtensionYAML    = ".yaml"
	ExtensionYML     = ".yml"
	ExtensionSQL     = ".sql"
	ExtensionProto   = ".proto"
)

var (
	// KnownExtensions is a list of all the known extensions.
	KnownExtensions = []string{
		ExtensionCSS,
		ExtensionGo,
		ExtensionHTML,
		ExtensionJS,
		ExtensionJSX,
		ExtensionPy,
		ExtensionSCSS,
		ExtensionSASS,
		ExtensionTS,
		ExtensionTSX,
		ExtensionYAML,
		ExtensionYML,
		ExtensionSQL,
		ExtensionProto,
	}

	// DefaultExtensionNoticeTemplates is a mapping between file extension (including the prefix dot) to the notice templates.
	DefaultExtensionNoticeTemplates = map[string]string{
		ExtensionCSS:   cssNoticeTemplate,
		ExtensionGo:    goNoticeTemplate,
		ExtensionHTML:  htmlNoticeTemplate,
		ExtensionJS:    jsNoticeTemplate,
		ExtensionJSX:   jsNoticeTemplate,
		ExtensionPy:    pythonNoticeTemplate,
		ExtensionSASS:  sassNoticeTemplate,
		ExtensionSCSS:  scssNoticeTemplate,
		ExtensionTS:    tsNoticeTemplate,
		ExtensionTSX:   tsNoticeTemplate,
		ExtensionYAML:  yamlNoticeTemplate,
		ExtensionYML:   yamlNoticeTemplate,
		ExtensionSQL:   sqlNoticeTemplate,
		ExtensionProto: protoNoticeTemplate,
	}

	// DefaultExcludes is the default excluded directories.
	DefaultExcludes = []string{
		".git/*",
		".github/*",
		"*/_config",
		"*/_config/*",
		"*/dist/*",
		"*/node_modules/*",
		"*/testdata",
		"*/testdata/*",
		"*/vendor/*",
		"node_modules/*",
		"protogen/*",
		"*.pb.go",
		"vendor/*",
		"venv/*",
		"*/venv/*",
		"*.vault.yml",
	}

	// DefaultIncludeFiles is the default included files list.
	DefaultIncludeFiles = []string{
		"*.css",
		"*.go",
		"*.html",
		"*.js",
		"*.jsx",
		"*.py",
		"*.sass",
		"*.scss",
		"*.ts",
		"*.tsx",
		"*.yaml",
		"*.yml",
		"*.sql",
		"*.proto",
	}
)

// Error Strings
var (
	VerifyErrorFormat = "%s: copyright header missing or invalid"
)

// Error sentinels
var (
	ErrWalkSkip = errors.New("walk skip; we should not process this file or path")
	ErrFailure  = errors.New("failure; one or more steps failed")
)

const (
	// goNoticeTemplate is the notice template specific to go files
	// note: it _must_ end in two newlines to prevent linting / compiler failures.
	goNoticeTemplate = `/*

{{ .Notice }}

*/

`

	yamlNoticeTemplate = `#
{{ .Notice | prefix "# " }}
#
`

	htmlNoticeTemplate = `<!--
{{ .Notice }}
-->
`

	jsNoticeTemplate = `/**
{{ .Notice | prefix " * " }}
 */
`

	tsNoticeTemplate = `/**
{{ .Notice | prefix " * " }}
 */
`

	cssNoticeTemplate = `/*
{{ .Notice | prefix " * " }}
 */
`

	scssNoticeTemplate = `/*
{{ .Notice | prefix " * " }}
 */
`

	sassNoticeTemplate = `/*
{{ .Notice | prefix " * " }}
 */
`

	pythonNoticeTemplate = `#
{{ .Notice | prefix "# " }}
#

`

	sqlNoticeTemplate = `--
{{ .Notice | prefix "-- " }}
--
`

	protoNoticeTemplate = `//
{{ .Notice | prefix "// " }}
//

`
)

const (
	goBuildTagExpr      = `^(\/\/(go:build| \+build).*\n)+\n`
	tsReferenceTagsExpr = `^(\/\/\/ \<reference path=\"(.*)\" \/\>\n)+`
	yearExpr            = `([0-9]{4,}?)`
	shebangExpr         = `(?s)^(\s*)#!([^\n]+)\n`
)

var (
	goBuildTagMatch      = regexp.MustCompile(goBuildTagExpr)
	tsReferenceTagsMatch = regexp.MustCompile(tsReferenceTagsExpr)
	yearMatch            = regexp.MustCompile(yearExpr)
	shebangMatch         = regexp.MustCompile(shebangExpr)
)
