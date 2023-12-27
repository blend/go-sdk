/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

// ConferenceBridge is an api type.
type ConferenceBridge struct {
	ConferenceNumber	string	`json:"conference_number,omitempty"`
	ConferenceURL		string	`json:"conference_url,omitempty"`
}
