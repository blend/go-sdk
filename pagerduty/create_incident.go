/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/r2"
)

// CreateIncident creates an incident using the incident events api.
func (hc HTTPClient) CreateIncident(ctx context.Context, incident CreateIncidentInput) (output Incident, err error) {
	var res *http.Response
	res, err = hc.Request(ctx,
		r2.OptPost(),
		r2.OptPath("/incidents"),
		r2.OptJSONBody(createIncidentInputWrapper{Incident: incident}),
	).Do()
	if err != nil {
		return
	}
	if statusCode := res.StatusCode; statusCode < 200 || statusCode > 299 {
		err = ex.New(ErrNon200Status, ex.OptMessagef("method: post, path: /incidents, status: %d", statusCode))
		return
	}
	defer res.Body.Close()
	var body createIncidentOutputWrapper
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		err = ex.New(err)
		return
	}
	output = body.Incident
	return
}

// CreateIncidentInput is the input to create|update incident.
type CreateIncidentInput struct {
	Type			string			`json:"type"`		// required
	Title			string			`json:"title"`		// required
	Service			APIObject		`json:"service"`	/// required
	Priority		*APIObject		`json:"priority,omitempty"`
	Body			*Body			`json:"body,omitempty"`
	IncidentKey		string			`json:"incident_key,omitempty"`
	Assignments		[]Assignment		`json:"assignments,omitempty"`
	EscalationPolicy	*APIObject		`json:"escalation_policy,omitempty"`
	Urgency			Urgency			`json:"urgency,omitempty"`
	ConferenceBridge	*ConferenceBridge	`json:"conference_bridge,omitempty"`
}

// createIncidentInputWrapper wraps the input to satisfy the input schema.
type createIncidentInputWrapper struct {
	Incident CreateIncidentInput `json:"incident"`
}

// CreateIncidentOutput is the response to create incident.
type createIncidentOutputWrapper struct {
	Incident Incident `json:"incident"`
}
