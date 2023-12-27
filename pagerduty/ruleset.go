/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package pagerduty

// RuleConditions represents the conditions field for a Ruleset
type RuleConditions struct {
	Operator          string              `json:"operator,omitempty"`
	RuleSubconditions []*RuleSubcondition `json:"subconditions,omitempty"`
}

// RuleSubcondition represents a subcondition of a ruleset condition
type RuleSubcondition struct {
	Operator   string              `json:"operator,omitempty"`
	Parameters *ConditionParameter `json:"parameters,omitempty"`
}

// ConditionParameter represents  parameters in a rule condition
type ConditionParameter struct {
	Path  string `json:"path,omitempty"`
	Value string `json:"value,omitempty"`
}

// RuleTimeFrame represents a time_frame object on the rule object
type RuleTimeFrame struct {
	ScheduledWeekly *ScheduledWeekly `json:"scheduled_weekly,omitempty"`
	ActiveBetween   *ActiveBetween   `json:"active_between,omitempty"`
}

// ScheduledWeekly represents a time_frame object for scheduling rules weekly
type ScheduledWeekly struct {
	Weekdays  []int  `json:"weekdays,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
	StartTime int    `json:"start_time,omitempty"`
	Duration  int    `json:"duration,omitempty"`
}

// ActiveBetween represents an active_between object for setting a timeline for rules
type ActiveBetween struct {
	StartTime int `json:"start_time,omitempty"`
	EndTime   int `json:"end_time,omitempty"`
}

// RuleActionParameter represents a generic parameter object on a rule action
type RuleActionParameter struct {
	Value string `json:"value,omitempty"`
}

// RuleActionSuppress represents a rule suppress action object
type RuleActionSuppress struct {
	Value               bool   `json:"value,omitempty"`
	ThresholdValue      int    `json:"threshold_value,omitempty"`
	ThresholdTimeUnit   string `json:"threshold_time_unit,omitempty"`
	ThresholdTimeAmount int    `json:"threshold_time_amount,omitempty"`
}

// RuleActionSuspend represents a rule suspend action object
type RuleActionSuspend struct {
	Value *bool `json:"value,omitempty"`
}

// RuleActionExtraction represents a rule extraction action object
type RuleActionExtraction struct {
	Target string `json:"target,omitempty"`
	Source string `json:"source,omitempty"`
	Regex  string `json:"regex,omitempty"`
}
