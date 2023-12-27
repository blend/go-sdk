/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*
Package configmeta provides a configutil metadata type to provide a canonical location to hold common config variables.

It provides a couple common variables to set with ldflags on build, namely `configmeta.Version` and `configmeta.GitRef`.

These can be set at build time with `go install -ldflags="-X github.com/blend/go-sdk/configmeta.Version=$(cat ${REPO_ROOT}/VERSION)" project/myapp` as an example.

The typical usage for the configmeta.Meta type is to embed in a config type and resolve it in your resolver.

Config Example:

	    type Config struct {
			configmeta.Meta `yaml:",inline"`
		}

		// Resolve resolves the config.
		func (c *Config) Resolve(ctx context.Context) error {
			return configutil.Resolve(ctx,
				(&c.Meta).Resolve,
			)
		}

This will pull `SERVICE_NAME` and `SERVICE_ENV` into relevant fields, as well as `configmeta.Version` into the Version field.

This type is used in a number of other packages for common fields like service name and service environment.
*/
package configmeta // import "github.com/blend/go-sdk/configmeta"
