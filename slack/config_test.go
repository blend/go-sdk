package slack

import "github.com/blend/go-sdk/configutil"

var (
	_ configutil.Resolver = (*Config)(nil)
)
