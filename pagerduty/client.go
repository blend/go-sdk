/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
