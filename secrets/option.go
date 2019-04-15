package secrets

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/blend/go-sdk/ex"
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

// OptRootCAs sets the root ca pool for client requests.
// If unset, it will set the VaultClient Client to be an http.Client.
// If unset, it will set the transport to be an http.Transport.
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
			var client *http.Client
			if vc.Client == nil {
				client = &http.Client{}
				vc.Client = client
			} else if typed, ok := vc.Client.(*http.Client); ok && typed != nil {
				client = typed
			}
			if client == nil {
				return ex.New("invalid http client for vault client; cannot set root cas")
			}

			var xport *http.Transport
			if client.Transport == nil {
				xport := &http.Transport{}
				client.Transport = xport
			} else if typed, ok := client.Transport.(*http.Transport); ok && typed != nil {
				xport = typed
			}
			if xport == nil {
				return ex.New("invalid http transport for vault client; cannot set root cas")
			}
			if xport.TLSClientConfig == nil {
				xport.TLSClientConfig = &tls.Config{}
			}
			xport.TLSClientConfig.RootCAs = certPool.Pool()
		}
		return nil
	}
}
