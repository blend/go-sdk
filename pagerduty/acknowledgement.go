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
