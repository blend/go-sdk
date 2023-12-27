/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

// ResolveReason is an api type.
type ResolveReason struct {
	Type		string		`json:"type,omitempty"`
	Incident	APIObject	`json:"incident,omitempty"`
}
