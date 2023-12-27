/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redis

import (
	"time"

	"github.com/blend/go-sdk/logger"
)

// OptConfig sets the redis config.
//
// Note: this will overwrite any existing settings.
func OptConfig(cfg Config) Option {
	return func(rc *RadixClient) error {
		rc.Config = cfg
		return nil
	}
}

// OptNetwork sets the redis network.
func OptNetwork(network string) Option {
	return func(rc *RadixClient) error {
		rc.Config.Network = network
		return nil
	}
}

// OptAddr sets the redis address.
func OptAddr(addr string) Option {
	return func(rc *RadixClient) error {
		rc.Config.Addr = addr
		return nil
	}
}

// OptAuthUser sets the redis auth user.
func OptAuthUser(user string) Option {
	return func(rc *RadixClient) error {
		rc.Config.AuthUser = user
		return nil
	}
}

// OptAuthPassword sets the redis auth password.
func OptAuthPassword(password string) Option {
	return func(rc *RadixClient) error {
		rc.Config.AuthPassword = password
		return nil
	}
}

// OptDB sets the redis db.
func OptDB(db string) Option {
	return func(rc *RadixClient) error {
		rc.Config.DB = db
		return nil
	}
}

// OptConnectTimeout sets the redis connect timeout.
func OptConnectTimeout(connectTimeout time.Duration) Option {
	return func(rc *RadixClient) error {
		rc.Config.ConnectTimeout = connectTimeout
		return nil
	}
}

// OptTimeout sets the redis general timeout.
func OptTimeout(timeout time.Duration) Option {
	return func(rc *RadixClient) error {
		rc.Config.Timeout = timeout
		return nil
	}
}

// OptLog sets the logger.
func OptLog(log logger.Triggerable) Option {
	return func(rc *RadixClient) error {
		rc.Log = log
		return nil
	}
}

// OptTracer sets the tracer.
func OptTracer(tracer Tracer) Option {
	return func(rc *RadixClient) error {
		rc.Tracer = tracer
		return nil
	}
}

// Option mutates a radix client.
type Option func(*RadixClient) error
