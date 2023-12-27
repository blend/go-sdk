/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

// APIObject represents generic api json response that is shared by most
// domain object
type APIObject struct {
	ID	string		`json:"id"`
	Type	ReferenceType	`json:"type"`
	Summary	string		`json:"summary,omitempty"`
	Self	string		`json:"self,omitempty"`
	HTMLUrl	string		`json:"html_url,omitempty"`
}

// APIReference are the fields required to reference another API object.
type APIReference struct {
	ID	string	`json:"id,omitempty"`
	Type	string	`json:"type,omitempty"`
}
