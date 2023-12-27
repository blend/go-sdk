/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package datadog

import "github.com/blend/go-sdk/configutil"

var (
	_ configutil.Resolver = (*Config)(nil)
)
