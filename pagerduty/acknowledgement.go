/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

import "time"

// Acknowledgement is an api type.
type Acknowledgement struct {
	At		time.Time	`json:"at,omitempty"`
	Acknowledger	APIObject	`json:"acknowledger,omitempty"`
}
