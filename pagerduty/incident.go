/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

import "time"

// Incident is the full api type for incidents.
type Incident struct {
	ID			string			`json:"id"`
	Summary			string			`json:"summary,omitempty"`
	Type			string			`json:"type,omitempty"`
	Self			string			`json:"self,omitempty"`
	HTMLUrl			string			`json:"html_url,omitempty"`
	IncidentNumber		int			`json:"incident_number,omitempty"`
	CreatedAt		time.Time		`json:"created_at,omitempty"`
	Status			IncidentStatus		`json:"status"`
	Title			string			`json:"title,omitempty"`
	PendingActions		[]Action		`json:"pending_actions,omitempty"`
	IncidentKey		string			`json:"incident_key,omitempty"`
	Service			APIObject		`json:"service,omitempty"`
	Assignments		[]Assignment		`json:"assignments,omitempty"`
	AssignedVia		string			`json:"assigned_via,omitempty"`
	Acknowledgements	[]Acknowledgement	`json:"acknowledgements,omitempty"`
	LastStatusChangeAt	time.Time		`json:"last_status_change_at,omitempty"`
	LastStatusChangeBy	APIObject		`json:"last_status_change_by,omitempty"`
	EscalationPolicy	APIObject		`json:"escalation_policy,omitempty"`
	Teams			[]APIObject		`json:"teams,omitempty"`
	Priority		APIObject		`json:"priority,omitempty"`
	Urgency			string			`json:"urgency"`
	ResolveReason		ResolveReason		`json:"resolve_reason,omitempty"`
	AlertCounts		struct {
		Triggered	int	`json:"triggered,omitempty"`
		Resolved	int	`json:"resolved,omitempty"`
		All		int	`json:"all,omitempty"`
	}	`json:"alert_counts,omitempty"`
	Body	Body	`json:"body,omitempty"`
}
