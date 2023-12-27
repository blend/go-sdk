/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

import "time"

// Action is an api type.
type Action struct {
	Type	string		`json:"type,omitempty"`
	At	time.Time	`json:"at,omitempty"`
}
