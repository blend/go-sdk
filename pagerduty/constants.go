/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

const (
	// DefaultAddr is the default api address.
	DefaultAddr = "https://api.pagerduty.com"
)

// ReferenceType is a type of reference.
type ReferenceType string

// ReferenceTypes
const (
	ReferenceTypeEscalationPolicy	ReferenceType	= "escalation_policy_reference"
	ReferenceTypeService		ReferenceType	= "service_reference"
	ReferenceTypeUser		ReferenceType	= "user_reference"
)

// Include is an object type constant.
type Include string

// Includes
const (
	IncludeUsers			Include	= "users"
	IncludeServices			Include	= "services"
	IncludeFirstTriggerLogEntries	Include	= "first_trigger_log_entries"
	IncludeEscalationPolicies	Include	= "escalation_policies"
	IncludeTeams			Include	= "teams"
	IncludeAssignees		Include	= "assignees"
	IncludeAcknowledgers		Include	= "acknowledgers"
	IncludePriorities		Include	= "priorities"
	IncludeConferenceBridge		Include	= "conference_bridge"
)

// IncidentStatus is a status for an incident
type IncidentStatus string

// IncidentStatuses
const (
	IncidentStatusTriggered		IncidentStatus	= "triggered"
	IncidentStatusAcknowledged	IncidentStatus	= "acknowledged"
	IncidentStatusResolved		IncidentStatus	= "resolved"
)

// Urgency is a urgency.
type Urgency string

// Urgencies
const (
	UrgencyHigh	Urgency	= "high"
	UrgencyLow	Urgency	= "low"
)
