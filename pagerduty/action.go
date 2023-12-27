/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package pagerduty

import "time"

// Action is an api type.
type Action struct {
	Type	string		`json:"type,omitempty"`
	At	time.Time	`json:"at,omitempty"`
}
