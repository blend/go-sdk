/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package statsutil

import (
	"github.com/blend/go-sdk/configmeta"
	"github.com/blend/go-sdk/datadog"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/statsd"
)

// MultiCollectorOptions are the options for the multi-collector.
type MultiCollectorOptions struct {
	configmeta.Meta

	DefaultTags	[]string
	Datadog		datadog.Config
	Prometheus	statsd.Config
	Printer		bool
}

// MultiCollectorOption mutates MultiCollectorOptions.
type MultiCollectorOption func(*MultiCollectorOptions)

// OptServiceName sets the default service name.
func OptServiceName(serviceName string) MultiCollectorOption {
	return func(mco *MultiCollectorOptions) {
		mco.ServiceName = serviceName
	}
}

// OptServiceEnv sets the default service env.
func OptServiceEnv(serviceEnv string) MultiCollectorOption {
	return func(mco *MultiCollectorOptions) {
		mco.ServiceEnv = serviceEnv
	}
}

// OptVersion sets the default version tag.
func OptVersion(version string) MultiCollectorOption {
	return func(mco *MultiCollectorOptions) {
		mco.Version = version
	}
}

// OptHostname sets the default hostname.
func OptHostname(hostname string) MultiCollectorOption {
	return func(mco *MultiCollectorOptions) {
		mco.Hostname = hostname
	}
}

// OptDefaultTags adds default tags.
func OptDefaultTags(tags ...string) MultiCollectorOption {
	return func(mco *MultiCollectorOptions) {
		mco.DefaultTags = append(mco.DefaultTags, tags...)
	}
}

// OptMetaConfig sets the datadog config.
func OptMetaConfig(meta configmeta.Meta) MultiCollectorOption {
	return func(mco *MultiCollectorOptions) {
		mco.Meta = meta
	}
}

// OptDatadogConfig sets the datadog config.
func OptDatadogConfig(cfg datadog.Config) MultiCollectorOption {
	return func(mco *MultiCollectorOptions) {
		mco.Datadog = cfg
	}
}

// OptPrometheusConfig sets the prometheus config.
func OptPrometheusConfig(cfg statsd.Config) MultiCollectorOption {
	return func(mco *MultiCollectorOptions) {
		mco.Prometheus = cfg
	}
}

// OptPrinter sets if we should enable the printer.
func OptPrinter(printer bool) MultiCollectorOption {
	return func(mco *MultiCollectorOptions) {
		mco.Printer = printer
	}
}

// NewMultiCollector initializes the stats collector(s).
func NewMultiCollector(log logger.Log, opts ...MultiCollectorOption) (stats.Collector, error) {
	var options MultiCollectorOptions
	for _, opt := range opts {
		opt(&options)
	}

	var collector stats.MultiCollector
	if options.Printer {
		logger.MaybeInfof(log, "using debug statsd printer")
		collector = append(collector, stats.NewPrinter(log))
	}

	if !options.Datadog.IsZero() {
		logger.MaybeInfof(log, "using datadog statsd collector: %s", options.Datadog.GetAddress())
		dd, err := datadog.New(options.Datadog)
		if err != nil {
			return nil, err
		}
		collector = append(collector, dd)
	} else {
		logger.MaybeDebugf(log, "datadog config unset, skipping")
	}

	if !options.Prometheus.IsZero() {
		logger.MaybeInfof(log, "using prometheus stats collector: %s", options.Prometheus.Addr)
		statsd, err := statsd.New(statsd.OptConfig(options.Prometheus))
		if err != nil {
			return nil, err
		}
		collector = append(collector, statsd)
	} else {
		logger.MaybeDebugf(log, "prometheus config unset, skipping")
	}

	// add default tags if there are collectors provisioned
	if len(collector) > 0 {
		if options.Meta.ServiceName != "" {
			collector.AddDefaultTags(stats.Tag(stats.TagService, options.Meta.ServiceName))
		} else if env.Env().ServiceName() != "" {
			collector.AddDefaultTags(stats.Tag(stats.TagService, env.Env().ServiceName()))
		}
		if options.Meta.ServiceEnv != "" {
			collector.AddDefaultTags(stats.Tag(stats.TagEnv, options.Meta.ServiceEnv))
		} else if env.Env().ServiceEnv() != "" {
			collector.AddDefaultTags(stats.Tag(stats.TagEnv, env.Env().ServiceEnv()))
		}
		if options.Meta.Version != "" {
			collector.AddDefaultTags(stats.Tag(stats.TagVersion, options.Meta.Version))
		} else if env.Env().Version() != "" {
			collector.AddDefaultTags(stats.Tag(stats.TagVersion, env.Env().Version()))
		}
		if options.Meta.ClusterName != "" {
			collector.AddDefaultTags(stats.Tag(stats.TagClusterName, options.Meta.ClusterName))
		} else if env.Env().ClusterName() != "" {
			collector.AddDefaultTags(stats.Tag(stats.TagClusterName, env.Env().ClusterName()))
		}
		collector.AddDefaultTags(options.DefaultTags...)
	}
	return collector, nil
}
