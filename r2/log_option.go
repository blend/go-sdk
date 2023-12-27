/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import "github.com/blend/go-sdk/sanitize"

// LogOptions are options that govern the logging of requests.
type LogOptions struct {
	RequestSanitizationDefaults []sanitize.RequestOption
}

// LogOption are mutators for log options.
type LogOption func(*LogOptions)
