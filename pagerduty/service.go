/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package pagerduty

// Integration is an endpoint (like Nagios, email, or an API call) that generates events, which are normalized and de-duplicated by PagerDuty to create incidents.
type Integration struct {
	APIObject
	Name             string     `json:"name,omitempty"`
	Service          *APIObject `json:"service,omitempty"`
	CreatedAt        string     `json:"created_at,omitempty"`
	Vendor           *APIObject `json:"vendor,omitempty"`
	Type             string     `json:"type,omitempty"`
	IntegrationKey   string     `json:"integration_key,omitempty"`
	IntegrationEmail string     `json:"integration_email,omitempty"`
}

// InlineModel represents when a scheduled action will occur.
type InlineModel struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

// ScheduledAction contains scheduled actions for the service.
type ScheduledAction struct {
	Type      string      `json:"type,omitempty"`
	At        InlineModel `json:"at,omitempty"`
	ToUrgency string      `json:"to_urgency"`
}

// IncidentUrgencyType are the incidents urgency during or outside support hours.
type IncidentUrgencyType struct {
	Type    string `json:"type,omitempty"`
	Urgency string `json:"urgency,omitempty"`
}

// SupportHours are the support hours for the service.
type SupportHours struct {
	Type       string `json:"type,omitempty"`
	Timezone   string `json:"time_zone,omitempty"`
	StartTime  string `json:"start_time,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
	DaysOfWeek []uint `json:"days_of_week,omitempty"`
}

// IncidentUrgencyRule is the default urgency for new incidents.
type IncidentUrgencyRule struct {
	Type                string               `json:"type,omitempty"`
	Urgency             string               `json:"urgency,omitempty"`
	DuringSupportHours  *IncidentUrgencyType `json:"during_support_hours,omitempty"`
	OutsideSupportHours *IncidentUrgencyType `json:"outside_support_hours,omitempty"`
}

// ListServiceRulesResponse represents a list of rules in a service
type ListServiceRulesResponse struct {
	Offset uint          `json:"offset,omitempty"`
	Limit  uint          `json:"limit,omitempty"`
	More   bool          `json:"more,omitempty"`
	Total  uint          `json:"total,omitempty"`
	Rules  []ServiceRule `json:"rules,omitempty"`
}

// ServiceRule represents a Service rule
type ServiceRule struct {
	ID         string              `json:"id,omitempty"`
	Self       string              `json:"self,omitempty"`
	Disabled   *bool               `json:"disabled,omitempty"`
	Conditions *RuleConditions     `json:"conditions,omitempty"`
	TimeFrame  *RuleTimeFrame      `json:"time_frame,omitempty"`
	Position   *int                `json:"position,omitempty"`
	Actions    *ServiceRuleActions `json:"actions,omitempty"`
}

// ServiceRuleActions represents a rule action
type ServiceRuleActions struct {
	Annotate    *RuleActionParameter   `json:"annotate,omitempty"`
	EventAction *RuleActionParameter   `json:"event_action,omitempty"`
	Extractions []RuleActionExtraction `json:"extractions,omitempty"`
	Priority    *RuleActionParameter   `json:"priority,omitempty"`
	Severity    *RuleActionParameter   `json:"severity,omitempty"`
	Suppress    *RuleActionSuppress    `json:"suppress,omitempty"`
	Suspend     *RuleActionSuspend     `json:"suspend,omitempty"`
}

// Service represents something you monitor (like a web service, email service, or database service).
type Service struct {
	APIObject
	Name                    string                   `json:"name,omitempty"`
	Description             string                   `json:"description,omitempty"`
	AutoResolveTimeout      *uint                    `json:"auto_resolve_timeout"`
	AcknowledgementTimeout  *uint                    `json:"acknowledgement_timeout"`
	CreateAt                string                   `json:"created_at,omitempty"`
	Status                  string                   `json:"status,omitempty"`
	LastIncidentTimestamp   string                   `json:"last_incident_timestamp,omitempty"`
	Integrations            []Integration            `json:"integrations,omitempty"`
	EscalationPolicy        EscalationPolicy         `json:"escalation_policy,omitempty"`
	Teams                   []Team                   `json:"teams,omitempty"`
	IncidentUrgencyRule     *IncidentUrgencyRule     `json:"incident_urgency_rule,omitempty"`
	SupportHours            *SupportHours            `json:"support_hours"`
	ScheduledActions        []ScheduledAction        `json:"scheduled_actions"`
	AlertCreation           string                   `json:"alert_creation,omitempty"`
	AlertGrouping           string                   `json:"alert_grouping,omitempty"`
	AlertGroupingTimeout    *uint                    `json:"alert_grouping_timeout,omitempty"`
	AlertGroupingParameters *AlertGroupingParameters `json:"alert_grouping_parameters,omitempty"`
}

// AlertGroupingParameters defines how alerts on the servicewill be automatically grouped into incidents
type AlertGroupingParameters struct {
	Type   string                 `json:"type"`
	Config AlertGroupParamsConfig `json:"config"`
}

// AlertGroupParamsConfig is the config object on alert_grouping_parameters
type AlertGroupParamsConfig struct {
	Timeout   uint     `json:"timeout,omitempty"`
	Aggregate string   `json:"aggregate,omitempty"`
	Fields    []string `json:"fields,omitempty"`
}
