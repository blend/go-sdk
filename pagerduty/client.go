/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

import "context"

// Client is the interface pagerduty clients implement.
type Client interface {
	CreateIncident(context.Context, CreateIncidentInput) (Incident, error)
	UpdateIncident(context.Context, string, UpdateIncidentInput) (Incident, error)
	ListIncidents(context.Context, ...ListIncidentOption) (ListIncidentsOutput, error)
	GetService(context.Context, string) (Service, error)
}
