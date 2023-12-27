/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package pagerduty

// Body is an api type.
type Body struct {
	Type	string	`json:"type"`
	Details	string	`json:"details,omitempty"`
}
