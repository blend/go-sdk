/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package copyright

// Option is a function that modifies a config.
type Option func(*Copyright)

// OptVerbose sets if we should show verbose output.
func OptVerbose(verbose bool) Option {
	return func(p *Copyright) {
		p.Config.Verbose = &verbose
	}
}

// OptDebug sets if we should show debug output.
func OptDebug(debug bool) Option {
	return func(p *Copyright) {
		p.Config.Debug = &debug
	}
}

// OptExitFirst sets if we should stop after the first failure.
func OptExitFirst(exitFirst bool) Option {
	return func(p *Copyright) {
		p.Config.ExitFirst = &exitFirst
	}
}

// OptExcludes sets the exclude glob filters.
func OptExcludes(excludeGlobs ...string) Option {
	return func(p *Copyright) {
		p.Config.Excludes = excludeGlobs
	}
}

// OptIncludeFiles sets the include file glob filters.
func OptIncludeFiles(includeGlobs ...string) Option {
	return func(p *Copyright) {
		p.Config.IncludeFiles = includeGlobs
	}
}

// OptNoticeBodyTemplate sets the notice body template.
func OptNoticeBodyTemplate(noticeBodyTemplate string) Option {
	return func(p *Copyright) {
		p.Config.NoticeBodyTemplate = noticeBodyTemplate
	}
}

// OptYear sets the template year.
func OptYear(year int) Option {
	return func(p *Copyright) {
		p.Config.Year = year
	}
}

// OptCompany sets the template company.
func OptCompany(company string) Option {
	return func(p *Copyright) {
		p.Config.Company = company
	}
}

// OptLicense sets the template license.
func OptLicense(license string) Option {
	return func(p *Copyright) {
		p.Config.License = license
	}
}

// OptRestrictions sets the template restrictions.
func OptRestrictions(restrictions string) Option {
	return func(p *Copyright) {
		p.Config.Restrictions = restrictions
	}
}

// OptConfig sets the config in its entirety.
func OptConfig(cfg Config) Option {
	return func(p *Copyright) {
		p.Config = cfg
	}
}
