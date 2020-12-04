package oauth

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
		m.RedirectURL = cfg.RedirectURI
		m.HostedDomain = cfg.HostedDomain
		m.AllowedDomains = cfg.AllowedDomains
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
		m.RedirectURL = redirectURI
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

// OptAllowedDomains sets the manager allowedDomains.
func OptAllowedDomains(allowedDomains ...string) Option {
	return func(m *Manager) error {
		m.AllowedDomains = allowedDomains
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

// OptTracer sets the manager tracer.
func OptTracer(tracer Tracer) Option {
	return func(m *Manager) error {
		m.Tracer = tracer
		return nil
	}
}
