/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package copyright

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ref"
)

func Test_Config(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	var cfg Config

	its.Equal(DefaultNoticeBodyTemplate, cfg.NoticeBodyTemplateOrDefault())
	cfg.NoticeBodyTemplate = "not-" + DefaultNoticeBodyTemplate
	its.Equal("not-"+DefaultNoticeBodyTemplate, cfg.NoticeBodyTemplateOrDefault())

	its.Equal(time.Now().UTC().Year(), cfg.YearOrDefault())
	cfg.Year = time.Now().UTC().Year() - 10
	its.Equal(time.Now().UTC().Year()-10, cfg.YearOrDefault())

	its.Equal(DefaultCompany, cfg.CompanyOrDefault())
	cfg.Company = "not-" + DefaultCompany
	its.Equal("not-"+DefaultCompany, cfg.CompanyOrDefault())

	its.Equal(DefaultOpenSourceLicense, cfg.LicenseOrDefault())
	cfg.License = "not-" + DefaultOpenSourceLicense
	its.Equal("not-"+DefaultOpenSourceLicense, cfg.LicenseOrDefault())

	its.Equal(DefaultRestrictionsInternal, cfg.RestrictionsOrDefault())
	cfg.Restrictions = "not-" + DefaultRestrictionsInternal
	its.Equal("not-"+DefaultRestrictionsInternal, cfg.RestrictionsOrDefault())

	its.Equal(DefaultExtensionNoticeTemplates, cfg.ExtensionNoticeTemplatesOrDefault())
	cfg.ExtensionNoticeTemplates = map[string]string{"foo": "bar"}
	its.Equal("bar", cfg.ExtensionNoticeTemplatesOrDefault()["foo"])

	its.False(cfg.ExitFirstOrDefault())
	cfg.ExitFirst = ref.Bool(true)
	its.True(cfg.ExitFirstOrDefault())

	its.False(cfg.QuietOrDefault())
	cfg.Quiet = ref.Bool(true)
	its.True(cfg.QuietOrDefault())

	its.False(cfg.VerboseOrDefault())
	cfg.Verbose = ref.Bool(true)
	its.True(cfg.VerboseOrDefault())

	its.False(cfg.DebugOrDefault())
	cfg.Debug = ref.Bool(true)
	its.True(cfg.DebugOrDefault())

	its.True(cfg.ShowDiffOrDefault())
	cfg.ShowDiff = ref.Bool(false)
	its.False(cfg.ShowDiffOrDefault())
}
