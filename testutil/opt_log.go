/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package testutil

import "github.com/blend/go-sdk/logger"

// OptLog sets the suite logger.
func OptLog(log logger.Log) Option {
	return func(s *Suite) {
		s.Log = log
	}
}
