/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package copyright

import "time"

// Config holds the runtime configuration option for the copyright engine.
type Config struct {
	// NoticeBodyTemplate is the notice body template that will be processed and
	// injected to the relevant extension specific notice template.
	NoticeBodyTemplate	string	`yaml:"noticeBodyTemplate"`
	// Year is the year to insert into the templates as `{{ .Year }}`
	Year	int	`yaml:"year"`
	// Company is the company name to insert into the templates as `{{ .Company }}`
	Company	string	`yaml:"company"`
	// License is the open source license to insert into in templates as `{{ .License }}`
	License	string	`yaml:"openSourceLicense"`

	// Restrictions an optional template to clarify copyright restrictions or
	// visibility modifiers, which is available in the `NoticeBodyTemplate` as `{{ .Restrictions }}`
	Restrictions	string	`yaml:"restrictionTemplate"`

	// Excludes are a list of globs to exclude, they can
	// match both files and directories.
	// This can be populated with `.gitignore` and the like.
	Excludes	[]string	`yaml:"excludes"`
	// IncludeFiles are a list of globs to match files to include.
	IncludeFiles	[]string	`yaml:"includeFiles"`

	// ExtensionNoticeTemplates is a map between file extension (including dot prefix)
	// to the relevant full notice template for the file. It can include a template variable
	// {{ .Notice }} that will insert the compiled `NoticyBodyTemplate`.
	ExtensionNoticeTemplates	map[string]string

	// FallbackNoticeTemplate is a full notice template that will be used if there is no extension
	// specific notice template.
	// It can include the template variable {{ .Notice }} that will instert the compiled `NoticeBodyTemplate`.
	FallbackNoticeTemplate	string

	// ExitFirst indicates if we should return after the first failure.
	ExitFirst	*bool	`yaml:"exitFirst"`
	// Quiet controls whether output is suppressed.
	Quiet	*bool	`yaml:"quiet"`
	// Verbose controls whether verbose output is shown.
	Verbose	*bool	`yaml:"verbose"`
	// Debug controls whether debug output is shown.
	Debug	*bool	`yaml:"debug"`

	// ShowDiff shows shows the diffs on verification failues.
	ShowDiff	*bool	`yaml:"verifyDiff"`
}

// NoticeBodyTemplateOrDefault returns the notice body template or a default.
func (c Config) NoticeBodyTemplateOrDefault() string {
	if c.NoticeBodyTemplate != "" {
		return c.NoticeBodyTemplate
	}
	return DefaultNoticeBodyTemplate
}

// YearOrDefault returns the current year or a default.
func (c Config) YearOrDefault() int {
	if c.Year > 0 {
		return c.Year
	}
	return time.Now().UTC().Year()
}

// CompanyOrDefault returns a company name or a default.
func (c Config) CompanyOrDefault() string {
	if c.Company != "" {
		return c.Company
	}
	return DefaultCompany
}

// LicenseOrDefault returns an open source license or a default.
func (c Config) LicenseOrDefault() string {
	if c.License != "" {
		return c.License
	}
	return DefaultOpenSourceLicense
}

// RestrictionsOrDefault returns restrictions or a default.
func (c Config) RestrictionsOrDefault() string {
	if c.Restrictions != "" {
		return c.Restrictions
	}
	return DefaultRestrictionsInternal
}

// ExtensionNoticeTemplatesOrDefault returns mapping between file extensions (including dot) to
// the notice templates (i.e. how the template should be fully formatted per file type).
func (c Config) ExtensionNoticeTemplatesOrDefault() map[string]string {
	if c.ExtensionNoticeTemplates != nil {
		return c.ExtensionNoticeTemplates
	}
	return DefaultExtensionNoticeTemplates
}

// ExitFirstOrDefault returns a value or a default.
func (c Config) ExitFirstOrDefault() bool {
	if c.ExitFirst != nil {
		return *c.ExitFirst
	}
	return false
}

// QuietOrDefault returns a value or a default.
func (c Config) QuietOrDefault() bool {
	if c.Quiet != nil {
		return *c.Quiet
	}
	return false
}

// VerboseOrDefault returns a value or a default.
func (c Config) VerboseOrDefault() bool {
	if c.Verbose != nil {
		return *c.Verbose
	}
	return false
}

// DebugOrDefault returns a value or a default.
func (c Config) DebugOrDefault() bool {
	if c.Debug != nil {
		return *c.Debug
	}
	return false
}

// ShowDiffOrDefault returns a value or a default.
func (c Config) ShowDiffOrDefault() bool {
	if c.ShowDiff != nil {
		return *c.ShowDiff
	}
	return true
}
