/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

// Body is an api type.
type Body struct {
	Type	string	`json:"type"`
	Details	string	`json:"details,omitempty"`
}
