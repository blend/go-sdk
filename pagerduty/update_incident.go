/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package pagerduty

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/r2"
)

// UpdateIncident implements the update incident method for the client.
func (hc HTTPClient) UpdateIncident(ctx context.Context, id string, incident UpdateIncidentInput) (output Incident, err error) {
	var res *http.Response
	res, err = hc.Request(ctx,
		r2.OptPut(),
		r2.OptPathf("/incidents/%s", id),
		r2.OptJSONBody(updateIncidentInputWrapper{Incident: incident}),
	).Do()
	if err != nil {
		return
	}
	if statusCode := res.StatusCode; statusCode < 200 || statusCode > 299 {
		err = ex.New(ErrNon200Status, ex.OptMessagef("method: put, path: /incidents/%s, status: %d", id, statusCode))
		return
	}
	defer res.Body.Close()
	var body updateIncidentOutputWrapper
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		err = ex.New(err)
		return
	}
	output = body.Incident
	return
}

// UpdateIncidentInput is the input to update incident.
type UpdateIncidentInput struct {
	Type			string			`json:"type"`			// required
	Status			IncidentStatus		`json:"status,omitempty"`	// required
	Priority		*APIObject		`json:"priority,omitempty"`
	Resolution		string			`json:"resolution,omitempty"`
	Title			string			`json:"title,omitempty"`
	EscalationLevel		int			`json:"escalation_level,omitempty"`
	Assignments		[]Assignment		`json:"assignments,omitempty"`
	EscalationPolicy	*APIObject		`json:"escalation_policy,omitempty"`
	Urgency			string			`json:"urgency,omitempty"`
	ConferenceBridge	*ConferenceBridge	`json:"conference_bridge,omitempty"`
}

// updateIncidentInputWrapper wraps the input to satisfy the input schema.
type updateIncidentInputWrapper struct {
	Incident UpdateIncidentInput `json:"incident"`
}

// updateIncidentOutputWrapper is the response to update incident.
type updateIncidentOutputWrapper struct {
	Incident Incident `json:"incident"`
}
