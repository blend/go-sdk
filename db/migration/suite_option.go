/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package migration

import (
	"github.com/blend/go-sdk/logger"
)

// SuiteOption is an option for migration Suites
type SuiteOption func(s *Suite)

// OptGroups allows you to add groups to the Suite. If you want, multiple OptGroups can be applied to the same Suite.
// They are additive.
func OptGroups(groups ...*Group) SuiteOption {
	return func(s *Suite) {
		if len(s.Groups) == 0 {
			s.Groups = groups
		} else {
			s.Groups = append(s.Groups, groups...)
		}
	}
}

// OptLog allows you to add a logger to the Suite.
func OptLog(log logger.Log) SuiteOption {
	return func(s *Suite) {
		s.Log = log
	}
}
