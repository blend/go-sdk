package reverseproxy

// ServerConfig as config for the reverse proxy server.
type ServerConfig struct {
	BindAddr        string `json:"bindAddr" yaml:"bindAddr" env:"BIND_ADDR"`
	UpgradeBindAddr string `json:"upgradeBindAddr" yaml:"upgradeBindAddr" env:"UPGRADE_BIND_ADDR"`

	// TLS options
	TLSKeyPath  string   `json:"tlsKeyPath" yaml:"tlsKeyPath" env:"TLS_KEY_PATH"`
	TLSCertPath string   `json:"tlsCertPath" yaml:"tlsCertPath" env:"TLS_CERT_PATH"`
	TLSCAPaths  []string `json:"tlsCAPaths" yaml:"tlsCAPAths" env:"TLS_CA_PATHS,csv"`

	Upstream         string `json:"upstream" yaml:"upstream"`
	UseProxyProtocol bool   `json:"useProxyProtocol" yaml:"useProxyProtocol"`
}

// BindAddrOrDefault returns the bind addr or a default.
func (c ServerConfig) BindAddrOrDefault() string {
	if c.BindAddr != "" {
		return c.BindAddr
	}
	return ":8443"
}

// UpstreamOrDefault returns the upstream or a default.
func (c ServerConfig) UpstreamOrDefault() string {
	if c.Upstream != "" {
		return c.Upstream
	}
	return "127.0.0.1:5000"
}
