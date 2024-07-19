/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sanitize

// Default values for disallowed field names
// Note: the values are compared using `strings.EqualFold` so the casing shouldn't matter
var (
	DefaultSanitizationDisallowedHeaders     = []string{"authorization", "cookie", "set-cookie"}
	DefaultSanitizationDisallowedQueryParams = []string{"access_token", "client_secret"}
)
