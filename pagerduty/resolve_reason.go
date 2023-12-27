/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package pagerduty

// ResolveReason is an api type.
type ResolveReason struct {
	Type     string    `json:"type,omitempty"`
	Incident APIObject `json:"incident,omitempty"`
}
