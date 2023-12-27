/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

// EscalationRule is a rule for an escalation policy to trigger.
type EscalationRule struct {
	ID	string		`json:"id,omitempty"`
	Delay	uint		`json:"escalation_delay_in_minutes,omitempty"`
	Targets	[]APIObject	`json:"targets"`
}

// EscalationPolicy is a collection of escalation rules.
type EscalationPolicy struct {
	APIObject
	Name		string			`json:"name,omitempty"`
	EscalationRules	[]EscalationRule	`json:"escalation_rules,omitempty"`
	Services	[]APIObject		`json:"services,omitempty"`
	NumLoops	uint			`json:"num_loops,omitempty"`
	Teams		[]APIReference		`json:"teams"`
	Description	string			`json:"description,omitempty"`
	RepeatEnabled	bool			`json:"repeat_enabled,omitempty"`
}
