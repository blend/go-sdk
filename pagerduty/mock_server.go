/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/blend/go-sdk/webutil"
)

var (
	_ http.Handler = (*MockAPI)(nil)
)

// MockAPI implements methods that can be called with the client.
type MockAPI struct {
	ListIncidents	func() (ListIncidentsOutput, error)
	UpdateIncident	func(string, UpdateIncidentInput) (Incident, error)
	CreateIncident	func(CreateIncidentInput) (Incident, error)
}

// Handler implements http.Handler.
func (ma MockAPI) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		if req.URL.Path == "/incidents" {
			if ma.ListIncidents != nil {
				output, err := ma.ListIncidents()
				if err != nil {
					_ = webutil.WriteJSON(rw, http.StatusInternalServerError, err.Error())
					return
				}
				_ = webutil.WriteJSON(rw, http.StatusOK, output)
				return
			}
		}
	}
	if req.Method == http.MethodPut {
		if strings.HasPrefix(req.URL.Path, "/incidents/") {
			if ma.UpdateIncident != nil {
				incidentID := strings.TrimPrefix(req.URL.Path, "/incidents/")
				var incidentBody updateIncidentInputWrapper
				if err := json.NewDecoder(req.Body).Decode(&incidentBody); err != nil {
					http.Error(rw, err.Error(), http.StatusBadRequest)
					return
				}
				incident, err := ma.UpdateIncident(incidentID, incidentBody.Incident)
				if err != nil {
					_ = webutil.WriteJSON(rw, http.StatusInternalServerError, err.Error())
					return
				}
				_ = webutil.WriteJSON(rw, http.StatusOK, updateIncidentOutputWrapper{Incident: incident})
				return
			}
		}
	}
	if req.Method == http.MethodPost {
		if req.URL.Path == "/incidents" {
			var incidentBody createIncidentInputWrapper
			if err := json.NewDecoder(req.Body).Decode(&incidentBody); err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			}
			incident, err := ma.CreateIncident(incidentBody.Incident)
			if err != nil {
				_ = webutil.WriteJSON(rw, http.StatusInternalServerError, err.Error())
				return
			}
			_ = webutil.WriteJSON(rw, http.StatusOK, updateIncidentOutputWrapper{Incident: incident})
			return
		}
	}
	http.NotFound(rw, req)
}
