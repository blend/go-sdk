/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package statsd

import (
	"context"
	"time"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

var (
	_ configutil.Resolver = (*Config)(nil)
)

// Config is the set of options for the statsd client.
type Config struct {
	Addr		string			`json:"addr" yaml:"addr" env:"STATSD_ADDR"`
	DialTimeout	time.Duration		`json:"dialTimeout" yaml:"dialTimeout"`
	SampleRate	float64			`json:"sampleRate" yaml:"sampleRate"`
	MaxPacketSize	int			`json:"maxPacketSize" yaml:"maxPacketSize"`
	MaxBufferSize	int			`json:"maxBufferSize" yaml:"maxBufferSize"`
	DefaultTags	map[string]string	`json:"defaultTags" yaml:"defaultTags"`
}

// Resolve implements configutil.Resolver.
func (c *Config) Resolve(ctx context.Context) error {
	return env.GetVars(ctx).ReadInto(c)
}

// IsZero returns if the config is set or not.
func (c Config) IsZero() bool {
	return c.Addr == ""
}
