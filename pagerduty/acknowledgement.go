/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package pagerduty

import "time"

// Acknowledgement is an api type.
type Acknowledgement struct {
	At           time.Time `json:"at,omitempty"`
	Acknowledger APIObject `json:"acknowledger,omitempty"`
}
