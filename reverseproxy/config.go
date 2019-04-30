package reverseproxy

// ServerConfig as config for the reverse proxy server.
type ServerConfig struct {
	BindAddr        string `json:"bindAddr" yaml:"bindAddr" env:"BIND_ADDR"`
	UpgradeBindAddr string `json:"upgradeBindAddr" yaml:"upgradeBindAddr" env:"UPGRADE_BIND_ADDR"`

	// TLS options
	TLSKeyPath  string   `json:"tlsKeyPath" yaml:"tlsKeyPath" env:"TLS_KEY_PATH"`
	TLSCertPath string   `json:"tlsCertPath" yaml:"tlsCertPath" env:"TLS_CERT_PATH"`
	TLSCAPaths  []string `json:"tlsCAPaths" yaml:"tlsCAPAths" env:"TLS_CA_PATHS,csv"`

	UseProxyProtocol bool `json:"useProxyProtocol" yaml:"useProxyProtocol"`
}
