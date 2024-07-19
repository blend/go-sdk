/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
