/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configmeta

import (
	"context"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

// These are set with `-ldflags="-X` on `go install`
var (
	// Version is the current version.
	Version = ""
	// GitRef is the currently deployed git ref
	GitRef = "HEAD"
	// ServiceName is the name of the service
	ServiceName = ""
	// ProjectName is the name of the project the service belongs to
	ProjectName = ""
	// ClusterName is the name of the cluster the service is deployed to
	ClusterName = ""
	// Region is the region the service is deployed to
	Region = ""
)

// Meta is the cluster config meta.
type Meta struct {
	// Region is the aws region the service is deployed to.
	Region string `json:"region,omitempty" yaml:"region,omitempty"`
	// ServiceName is name of the service
	ServiceName string `json:"serviceName,omitempty" yaml:"serviceName,omitempty"`
	// ProjectName is the project name injected by Deployinator.
	ProjectName string `json:"projectName,omitempty" yaml:"projectName,omitempty"`
	// ClusterName is the name of the cluster the service is deployed to
	ClusterName string `json:"clusterName,omitempty" yaml:"clusterName,omitempty"`
	// Environment is the environment of the cluster (sandbox, prod etc.)
	ServiceEnv string `json:"serviceEnv,omitempty" yaml:"serviceEnv,omitempty"`
	// Hostname is the environment of the cluster (sandbox, prod etc.)
	Hostname string `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	// Version is the application version.
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	// GitRef is the git ref of the image.
	GitRef string `json:"gitRef,omitempty" yaml:"gitRef,omitempty"`
}

// SetFrom returns a resolve action to set this meta from a root meta.
func (m *Meta) SetFrom(other *Meta) configutil.ResolveAction {
	return func(_ context.Context) error {
		m.Region = other.Region
		m.ServiceName = other.ServiceName
		m.ProjectName = other.ProjectName
		m.ClusterName = other.ClusterName
		m.ServiceEnv = other.ServiceEnv
		m.Hostname = other.Hostname
		m.Version = other.Version
		m.GitRef = other.GitRef
		return nil
	}
}

// ApplyTo applies a given meta to another meta.
func (m *Meta) ApplyTo(other *Meta) configutil.ResolveAction {
	return func(_ context.Context) error {
		other.Region = m.Region
		other.ServiceName = m.ServiceName
		other.ProjectName = m.ProjectName
		other.ClusterName = m.ClusterName
		other.ServiceEnv = m.ServiceEnv
		other.Hostname = m.Hostname
		other.Version = m.Version
		other.GitRef = m.GitRef
		return nil
	}
}

// Resolve implements configutil.Resolver
func (m *Meta) Resolve(ctx context.Context) error {
	return configutil.Resolve(ctx,
		configutil.SetString(&m.Region, configutil.Env(env.VarRegion), configutil.String(m.Region), configutil.LazyString(&Region)),
		configutil.SetString(&m.ServiceName, configutil.Env(env.VarServiceName), configutil.String(m.ServiceName), configutil.LazyString(&ServiceName)),
		configutil.SetString(&m.ProjectName, configutil.Env(env.VarProjectName), configutil.String(m.ProjectName), configutil.LazyString(&ProjectName)),
		configutil.SetString(&m.ClusterName, configutil.Env(env.VarClusterName), configutil.String(m.ClusterName), configutil.LazyString(&ClusterName)),
		configutil.SetString(&m.ServiceEnv, configutil.Env(env.VarServiceEnv), configutil.String(m.ServiceEnv), configutil.String(env.ServiceEnvDev)),
		configutil.SetString(&m.Hostname, configutil.Env(env.VarHostname), configutil.String(m.Hostname)),
		configutil.SetString(&m.Version, configutil.Env(env.VarVersion), configutil.String(m.Version), configutil.LazyString(&Version), configutil.String(DefaultVersion)),
		configutil.SetString(&m.GitRef, configutil.Env(env.VarGitRef), configutil.String(m.GitRef), configutil.LazyString(&GitRef)),
	)
}

// RegionOrDefault returns the region or the default.
func (m Meta) RegionOrDefault() string {
	if m.Region != "" {
		return m.Region
	}
	return Region
}

// ServiceNameOrDefault returns the service name or the default.
func (m Meta) ServiceNameOrDefault() string {
	if m.ServiceName != "" {
		return m.ServiceName
	}
	return ServiceName
}

// ProjectNameOrDefault returns the project name or the default.
func (m Meta) ProjectNameOrDefault() string {
	if m.ProjectName != "" {
		return m.ProjectName
	}
	return ProjectName
}

// ClusterNameOrDefault returns the cluster name or the default.
func (m Meta) ClusterNameOrDefault() string {
	if m.ClusterName != "" {
		return m.ClusterName
	}
	return ClusterName
}

// ServiceEnvOrDefault returns the cluster environment.
func (m Meta) ServiceEnvOrDefault() string {
	if m.ServiceEnv != "" {
		return m.ServiceEnv
	}
	return env.DefaultServiceEnv
}

// VersionOrDefault returns a version or a default.
func (m Meta) VersionOrDefault() string {
	if m.Version != "" {
		return m.Version
	}
	if Version != "" {
		return Version
	}
	return DefaultVersion
}

// GitRefOrDefault returns a gitref or a default.
func (m Meta) GitRefOrDefault() string {
	if m.GitRef != "" {
		return m.GitRef
	}
	if GitRef != "" {
		return GitRef
	}
	return "HEAD"
}

// IsProdlike returns if the cluster meta environment is prodlike.
func (m Meta) IsProdlike() bool {
	return env.IsProdlike(m.ServiceEnvOrDefault())
}
