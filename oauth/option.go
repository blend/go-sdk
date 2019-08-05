package oauth

import "github.com/blend/go-sdk/r2"

// Option is an option for oauth managers.
type Option func(*Manager) error

// OptConfig sets a manager based on a config.
func OptConfig(cfg Config) Option {
	return func(m *Manager) error {
		secret, err := cfg.DecodeSecret()
		if err != nil {
			return err
		}
		m.Secret = secret
		m.RedirectURI = cfg.RedirectURI
		m.HostedDomain = cfg.HostedDomain
		m.Scopes = cfg.ScopesOrDefault()
		m.ClientID = cfg.ClientID
		m.ClientSecret = cfg.ClientSecret
		return nil
	}
}

// OptClientID sets the manager cliendID.
func OptClientID(cliendID string) Option {
	return func(m *Manager) error {
		m.ClientID = cliendID
		return nil
	}
}

// OptClientSecret sets the manager clientSecret.
func OptClientSecret(clientSecret string) Option {
	return func(m *Manager) error {
		m.ClientSecret = clientSecret
		return nil
	}
}

// OptSecret sets the manager secret.
func OptSecret(secret []byte) Option {
	return func(m *Manager) error {
		m.Secret = secret
		return nil
	}
}

// OptRedirectURI sets the manager redirectURI.
func OptRedirectURI(redirectURI string) Option {
	return func(m *Manager) error {
		m.RedirectURI = redirectURI
		return nil
	}
}

// OptHostedDomain sets the manager hostedDomain.
func OptHostedDomain(hostedDomain string) Option {
	return func(m *Manager) error {
		m.HostedDomain = hostedDomain
		return nil
	}
}

// OptScopes sets the manager scopes.
func OptScopes(scopes ...string) Option {
	return func(m *Manager) error {
		m.Scopes = scopes
		return nil
	}
}

// OptRequestDefaults sets the request defaults.
// These affect making the profile fetch call
// to google and do not affect the token exchange.
func OptRequestDefaults(defaults ...r2.Option) Option {
	return func(m *Manager) error {
		m.RequestDefaults = defaults
		return nil
	}
}

// OptTracer sets the manager tracer.
func OptTracer(tracer Tracer) Option {
	return func(m *Manager) error {
		m.Tracer = tracer
		return nil
	}
}
