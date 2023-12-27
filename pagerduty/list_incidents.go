/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/r2"
)

// ListIncidents lists incidents.
//
// Use the variadic options to set constraining query parameters to filter or sort the incidents.
func (hc HTTPClient) ListIncidents(ctx context.Context, opts ...ListIncidentOption) (output ListIncidentsOutput, err error) {
	var options ListIncidentsOptions
	for _, opt := range opts {
		opt(&options)
	}
	var res *http.Response
	res, err = hc.Request(ctx,
		append([]r2.Option{
			r2.OptGet(),
			r2.OptPath("/incidents"),
		}, options.Options()...)...,
	).Do()
	if err != nil {
		return
	}
	if statusCode := res.StatusCode; statusCode < 200 || statusCode > 299 {
		err = ex.New(ErrNon200Status, ex.OptMessagef("method: post, path: /incidents, status: %d", statusCode))
		return
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&output); err != nil {
		err = ex.New(err)
		return
	}
	return
}

// OptListIncidentsDateRange sets a field on the options.
func OptListIncidentsDateRange(dateRange string) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.DateRange = dateRange }
}

// OptListIncidentsIncidentKey sets a field on the options.
func OptListIncidentsIncidentKey(incidentKey string) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.IncidentKey = incidentKey }
}

// OptListIncidentsInclude sets the "include" query string value.
//
// Include sets if we should add additional data to the response for
// corresponding fields on the output object.
func OptListIncidentsInclude(include ...Include) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.Include = include }
}

// OptListIncidentsLimit sets a field on the options.
func OptListIncidentsLimit(limit int) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.Limit = limit }
}

// OptListIncidentsOffset sets a field on the options.
func OptListIncidentsOffset(offset int) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.Offset = offset }
}

// OptListIncidentsServiceIDs sets a field on the options.
func OptListIncidentsServiceIDs(serviceIDs ...string) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.ServiceIDs = serviceIDs }
}

// OptListIncidentsSince sets a field on the options.
func OptListIncidentsSince(since string) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.Since = since }
}

// OptListIncidentsSortBy sets a field on the options.
func OptListIncidentsSortBy(sortBy string) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.SortBy = sortBy }
}

// OptListIncidentsStatuses sets a field on the options.
func OptListIncidentsStatuses(statuses ...IncidentStatus) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.Statuses = statuses }
}

// OptListIncidentsTeamIDs sets a field on the options.
func OptListIncidentsTeamIDs(teamIDs ...string) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.TeamIDs = teamIDs }
}

// OptListIncidentsTimeZone sets a field on the options.
func OptListIncidentsTimeZone(timeZone string) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.TimeZone = timeZone }
}

// OptListIncidentsTotal sets a field on the options.
func OptListIncidentsTotal(total bool) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.Total = &total }
}

// OptListIncidentsUntil sets a field on the options.
func OptListIncidentsUntil(until string) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.Until = until }
}

// OptListIncidentsUrgencies sets a field on the options.
func OptListIncidentsUrgencies(urgencies ...string) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.Urgencies = urgencies }
}

// OptListIncidentsUserIDs sets a field on the options.
func OptListIncidentsUserIDs(userIDs ...string) ListIncidentOption {
	return func(lio *ListIncidentsOptions) { lio.UserIDs = userIDs }
}

// ListIncidentOption mutates the list incidents options.
type ListIncidentOption func(*ListIncidentsOptions)

// ListIncidentsOptions are all the options for a list incidents call.
type ListIncidentsOptions struct {
	DateRange	string
	IncidentKey	string
	Include		[]Include
	Limit		int
	Offset		int
	ServiceIDs	[]string
	Since		string
	SortBy		string
	Statuses	[]IncidentStatus
	TeamIDs		[]string
	TimeZone	string
	Total		*bool
	Until		string
	Urgencies	[]string
	UserIDs		[]string
}

// Options yields the r2 options for the options.
//
// _Allow myself to introduce ... myself_
func (lio ListIncidentsOptions) Options() (output []r2.Option) {
	if lio.DateRange != "" {
		output = append(output, r2.OptQueryValue("date_range", lio.DateRange))
	}
	if lio.IncidentKey != "" {
		output = append(output, r2.OptQueryValue("incident_key", lio.IncidentKey))
	}
	if len(lio.Include) > 0 {
		for _, include := range lio.Include {
			output = append(output, r2.OptQueryValueAdd("include[]", string(include)))
		}
	}
	if lio.Limit > 0 {
		output = append(output, r2.OptQueryValue("limit ", fmt.Sprint(lio.Limit)))
	}
	if lio.Offset > 0 {
		output = append(output, r2.OptQueryValue("offset", fmt.Sprint(lio.Offset)))
	}
	if len(lio.ServiceIDs) > 0 {
		for _, serviceID := range lio.ServiceIDs {
			output = append(output, r2.OptQueryValueAdd("service_ids[]", serviceID))
		}
	}
	if lio.Since != "" {
		output = append(output, r2.OptQueryValue("since", lio.Since))
	}
	if lio.SortBy != "" {
		output = append(output, r2.OptQueryValue("sort_by", lio.SortBy))
	}
	if len(lio.Statuses) > 0 {
		for _, status := range lio.Statuses {
			output = append(output, r2.OptQueryValueAdd("statuses[]", string(status)))
		}
	}
	if len(lio.TeamIDs) > 0 {
		for _, teamID := range lio.TeamIDs {
			output = append(output, r2.OptQueryValueAdd("team_ids[]", teamID))
		}
	}
	if lio.TimeZone != "" {
		output = append(output, r2.OptQueryValue("time_zone", lio.TimeZone))
	}
	if lio.Total != nil {
		output = append(output, r2.OptQueryValue("total", fmt.Sprint(*lio.Total)))
	}
	if lio.Until != "" {
		output = append(output, r2.OptQueryValue("until", lio.Until))
	}
	if len(lio.Urgencies) > 0 {
		for _, urgency := range lio.Urgencies {
			output = append(output, r2.OptQueryValueAdd("urgencies[]", urgency))
		}
	}
	if len(lio.UserIDs) > 0 {
		for _, userID := range lio.UserIDs {
			output = append(output, r2.OptQueryValueAdd("user_ids[]", userID))
		}
	}
	return
}

// ListIncidentsOutput is the output of a list incidents call.
type ListIncidentsOutput struct {
	Offset		int		`json:"offset"`
	Limit		int		`json:"limit"`
	More		bool		`json:"more"`
	Total		*int		`json:"total"`
	Incidents	[]Incident	`json:"incidents"`
}
