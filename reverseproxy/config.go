package reverseproxy

// Config as config for the reverse proxy.
type Config struct {
	BindAddr         string `json:"bindAddr" yaml:"bindAddr" env:"BIND_ADDR"`
	RedirectBindAddr string `json:"redirectBindAddr" yaml:"redirectBindAddr" env:"REDIRECT_BIND_ADDR"`

	// TLS options
	TLSKeyPath  string   `json:"tlsKeyPath" yaml:"tlsKeyPath" env:"TLS_KEY_PATH"`
	TLSCertPath string   `json:"tlsCertPath" yaml:"tlsCertPath" env:"TLS_CERT_PATH"`
	TLSCAPaths  []string `json:"tlsCAPaths" yaml:"tlsCAPAths" env:"TLS_CA_PATHS,csv"`

	// Upstream options
	Upstreams UpstreamConfig `json:"upstreams" yaml:"upstreams"`
}

// UpstreamConfig is a config for upstreams.
type UpstreamConfig struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}
