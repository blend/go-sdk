/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package breaker

import "time"

// Config is the breaker config.
type Config struct {
	HalfOpenMaxActions	int64		`json:"halfOpenMaxActions" yaml:"halfOpenMaxActions"`
	ClosedExpiryInterval	time.Duration	`json:"closedExpiryInterval" yaml:"closedExpiryInterval"`
	OpenExpiryInterval	time.Duration	`json:"openExpiryInterval" yaml:"openExpiryInterval"`
}
