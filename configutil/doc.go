/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

/*
Package configutil contains helpers for configuration.

The intent of this package is to allow you to compose configuration types found throughout `sdk/` and elsewhere into a single configuration type for a given program.

An example configuration type might be:

	type Config struct {
		DB db.Config `yaml:"db"`
		Logger logger.Config `yaml:"logger"`
		Web web.Config `yaml:"web"`
	}

On start, you may want to read from a yaml configuration file into this type, deserializing the configuration.

	var cfg Config
	paths, err := configutil.Read(&cfg)
	// `err` is the first error returned from reading files
	// `paths` are the paths read from to populate the config

If we want to also read from the environment, we can do so by adding a resolver, and using the `configutil.Resolve`, `configutil.SetXYZ` helpers.

	func (c *Config) Resolve(ctx context.Context) error {
		// NOTE: the pointer receiver is critical to modify the fields of the config
		return configutil.Resolve(ctx,
			(&c.DB).Resolve, // we use the (&foo) form here because we want to call the pointer receiver for `Resolve`
			(&c.Logger).Resolve,
			(&c.Web).Resolve,

			configutil.SetString(&c.Web.BindAddr, configutil.Env("BIND_ADDR"), configutil.String(c.Web.BindAddr)),
		)
	}

In the above, the environment variable `BIND_ADDR` takes precedence over the string value found in any configuration file(s).

Note, we also "resolve" each of the attached configs first, in case they also have environment variables they read from etc.
*/
package configutil	// import "github.com/blend/go-sdk/configutil"
