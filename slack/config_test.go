/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package slack

import "github.com/blend/go-sdk/configutil"

var (
	_ configutil.Resolver = (*Config)(nil)
)
