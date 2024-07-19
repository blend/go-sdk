/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package copyright

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Options(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	c := new(Copyright)

	its.Nil(c.Config.Verbose)
	OptVerbose(true)(c)
	its.True(*c.Config.Verbose)

	its.Nil(c.Config.Debug)
	OptDebug(true)(c)
	its.True(*c.Config.Debug)

	its.Nil(c.Config.ExitFirst)
	OptExitFirst(true)(c)
	its.True(*c.Config.ExitFirst)

	its.Empty(c.Config.IncludeFiles)
	OptIncludeFiles("opt-include-0", "opt-include-1")(c)
	its.Equal([]string{"opt-include-0", "opt-include-1"}, c.Config.IncludeFiles)

	its.Empty(c.Config.Excludes)
	OptExcludes("opt-exclude-0", "opt-exclude-1")(c)
	its.Equal([]string{"opt-exclude-0", "opt-exclude-1"}, c.Config.Excludes)

	its.Empty(c.Config.NoticeBodyTemplate)
	OptNoticeBodyTemplate("opt-notice-body-template")(c)
	its.Equal("opt-notice-body-template", c.Config.NoticeBodyTemplate)

	its.Zero(c.Config.Year)
	OptYear(2021)(c)
	its.Equal(2021, c.Config.Year)

	its.Empty(c.Config.Company)
	OptCompany("opt-company")(c)
	its.Equal("opt-company", c.Config.Company)

	its.Empty(c.Config.License)
	OptLicense("opt-license")(c)
	its.Equal("opt-license", c.Config.License)

	its.Empty(c.Config.Restrictions)
	OptRestrictions("opt-restrictions")(c)
	its.Equal("opt-restrictions", c.Config.Restrictions)

	OptConfig(Config{})(c)
	its.Empty(c.Config.Restrictions)
}
