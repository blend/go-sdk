/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

import "time"

// Assignment is an assignment.
type Assignment struct {
	At		time.Time	`json:"at,omitempty"`
	Assignee	*APIObject	`json:"assignee"`
}
