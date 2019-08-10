package secrets

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/blend/go-sdk/logger"
)

// Option is an option for a vault client.
type Option func(*VaultClient) error

// OptLog sets the logger on the vault client.
func OptLog(log logger.Log) Option {
	return func(vc *VaultClient) error {
		vc.Log = log
		return nil
	}
}

// OptConfigFromEnv sets the vault client from a given configuration read
// from the environment.
func OptConfigFromEnv() Option {
	return func(vc *VaultClient) error {
		cfg, err := NewConfigFromEnv()
		if err != nil {
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
	return func(vc *VaultClient) error {
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
	return func(vc *VaultClient) error {
		remote, err := url.Parse(addr)
		if err != nil {
			return err
		}
		vc.Remote = remote
		return nil
	}
}

// OptMount sets the vault client mount.
func OptMount(mount string) Option {
	return func(vc *VaultClient) error {
		vc.Mount = mount
		return nil
	}
}

// OptToken sets the vault client token.
func OptToken(token string) Option {
	return func(vc *VaultClient) error {
		vc.Token = token
		return nil
	}
}

// OptTimeout sets the timeout to vault
func OptTimeout(timeout time.Duration) Option{
	return func(vc *VaultClient) error {
		vc.Timeout = timeout
		return nil
	}
}

// OptRootCAs sets the root ca pool for client requests.
func OptRootCAs(rootCAs ...string) Option {
	return func(vc *VaultClient) error {
		if len(rootCAs) > 0 {
			certPool, err := NewCertPool()
			if err != nil {
				return err
			}
			err = certPool.AddPaths(rootCAs...)
			if err != nil {
				return err
			}

			xport := &http.Transport{}
			xport.TLSClientConfig = &tls.Config{}
			xport.TLSClientConfig.RootCAs = certPool.Pool()
			vc.Transport = xport
		}
		return nil
	}
}

// OptTracer allows you to configure a tracer on the vault client
func OptTracer(tracer Tracer) Option {
	return func(vc *VaultClient) error {
		vc.Tracer = tracer
		return nil
	}
}
