/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sanitize

// Default values for disallowed field names
// Note: the values are compared using `strings.EqualFold` so the casing shouldn't matter
var (
	DefaultSanitizationDisallowedHeaders		= []string{"authorization", "cookie", "set-cookie"}
	DefaultSanitizationDisallowedQueryParams	= []string{"access_token", "client_secret"}
)
