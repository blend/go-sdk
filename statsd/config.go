package statsd

import (
	"context"
	"time"

	"github.com/blend/go-sdk/configutil"
)

var (
	_ configutil.Resolver = (*Config)(nil)
)

// Config is the set of options for the statsd client.
type Config struct {
	Addr          string            `json:"addr" yaml:"addr" env:"STATSD_ADDR"`
	DialTimeout   time.Duration     `json:"dialTimeout" yaml:"dialTimeout"`
	SampleRate    float64           `json:"sampleRate" yaml:"sampleRate"`
	MaxPacketSize int               `json:"maxPacketSize" yaml:"maxPacketSize"`
	DefaultTags   map[string]string `json:"defaultTags" yaml:"defaultTags"`
}

// Resolve implements configutil.Resolver.
func (c *Config) Resolve(ctx context.Context) error {
	return configutil.GetEnvVars(ctx).ReadInto(c)
}
