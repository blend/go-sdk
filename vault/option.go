/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// Option is an option for a vault client.
type Option func(*APIClient) error

// OptLog sets the logger on the vault client.
func OptLog(log logger.Log) Option {
	return func(vc *APIClient) error {
		vc.Log = log
		return nil
	}
}

// OptConfigFromEnv sets the vault client from a given configuration read
// from the environment.
func OptConfigFromEnv() Option {
	return func(vc *APIClient) error {
		var cfg Config
		if err := (&cfg).Resolve(env.WithVars(context.Background(), env.Env())); err != nil {
			return err
		}
		if err := OptConfig(cfg)(vc); err != nil {
			return err
		}
		return nil
	}
}

// OptConfig sets the vault client from a given configuration.
func OptConfig(cfg Config) Option {
	return func(vc *APIClient) error {
		if err := OptRemote(cfg.AddrOrDefault())(vc); err != nil {
			return err
		}
		if err := OptMount(cfg.MountOrDefault())(vc); err != nil {
			return err
		}
		if err := OptToken(cfg.Token)(vc); err != nil {
			return err
		}
		if err := OptTimeout(cfg.TimeoutOrDefault())(vc); err != nil {
			return err
		}
		if err := OptRootCAs(cfg.RootCAs...)(vc); err != nil {
			return err
		}
		return nil
	}
}

// OptRemote sets the client remote.
func OptRemote(addr string) Option {
	return func(vc *APIClient) error {
		remote, err := url.Parse(addr)
		if err != nil {
			return err
		}
		vc.Remote = remote
		return nil
	}
}

// OptAddr is an alias to OptRemote.
func OptAddr(addr string) Option {
	return OptRemote(addr)
}

// OptMount sets the vault client mount.
func OptMount(mount string) Option {
	return func(vc *APIClient) error {
		vc.Mount = mount
		return nil
	}
}

// OptToken sets the vault client token.
func OptToken(token string) Option {
	return func(vc *APIClient) error {
		vc.Token = token
		return nil
	}
}

// OptTimeout sets the timeout to vault
func OptTimeout(timeout time.Duration) Option {
	return func(vc *APIClient) error {
		vc.Timeout = timeout
		return nil
	}
}

// OptRootCAs sets the root ca pool for client requests.
func OptRootCAs(rootCAs ...string) Option {
	return func(vc *APIClient) error {
		if len(rootCAs) > 0 {
			certPool, err := x509.SystemCertPool()
			if err != nil {
				return err
			}

			for _, caPath := range rootCAs {
				contents, err := os.ReadFile(caPath)
				if err != nil {
					return err
				}
				if ok := certPool.AppendCertsFromPEM(contents); !ok {
					return ex.New("Invalid Root CA")
				}
			}

			xport := new(http.Transport)
			xport.TLSClientConfig = new(tls.Config)
			xport.TLSClientConfig.RootCAs = certPool
			vc.Transport = xport
		}
		return nil
	}
}

// OptTracer allows you to configure a tracer on the vault client
func OptTracer(tracer Tracer) Option {
	return func(vc *APIClient) error {
		vc.Tracer = tracer
		return nil
	}
}
