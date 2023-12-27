/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package pagerduty

// Team is a collection of users and escalation policies that represent a group of people within an organization.
type Team struct {
	APIObject
	Name		string	`json:"name,omitempty"`
	Description	string	`json:"description,omitempty"`
}
