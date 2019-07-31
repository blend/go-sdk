package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/webutil"
)

// Config is an object used to set up a web app.
type Config struct {
	Port                      int32         `json:"port,omitempty" yaml:"port,omitempty" env:"PORT"`
	BindAddr                  string        `json:"bindAddr,omitempty" yaml:"bindAddr,omitempty" env:"BIND_ADDR"`
	BaseURL                   string        `json:"baseURL,omitempty" yaml:"baseURL,omitempty" env:"BASE_URL"`
	SkipRedirectTrailingSlash bool          `json:"skipRedirectTrailingSlash,omitempty" yaml:"skipRedirectTrailingSlash,omitempty"`
	HandleOptions             bool          `json:"handleOptions,omitempty" yaml:"handleOptions,omitempty"`
	HandleMethodNotAllowed    bool          `json:"handleMethodNotAllowed,omitempty" yaml:"handleMethodNotAllowed,omitempty"`
	DisablePanicRecovery      bool          `json:"disablePanicRecovery,omitempty" yaml:"disablePanicRecovery,omitempty"`
	SessionTimeout            time.Duration `json:"sessionTimeout,omitempty" yaml:"sessionTimeout,omitempty" env:"SESSION_TIMEOUT"`
	SessionTimeoutIsRelative  bool          `json:"sessionTimeoutIsRelative,omitempty" yaml:"sessionTimeoutIsRelative,omitempty" env:"SESSION_TIMEOUT_RELATIVE"`

	CookieSecure   *bool  `json:"cookieSecure,omitempty" yaml:"cookieSecure,omitempty" env:"COOKIE_SECURE"`
	CookieHTTPOnly *bool  `json:"cookieHTTPOnly,omitempty" yaml:"cookieHTTPOnly,omitempty" env:"COOKIE_HTTP_ONLY"`
	CookieSameSite string `json:"cookieSameSite,omitempty" yaml:"cookieSameSite,omitempty" env:"COOKIE_SAME_SITE"`
	CookieName     string `json:"cookieName,omitempty" yaml:"cookieName,omitempty" env:"COOKIE_NAME"`
	CookiePath     string `json:"cookiePath,omitempty" yaml:"cookiePath,omitempty" env:"COOKIE_PATH"`
	CookieDomain   string `json:"cookieDomain,omitempty" yaml:"cookieDomain,omitempty" env:"COOKIE_DOMAIN"`

	DefaultHeaders      map[string]string `json:"defaultHeaders,omitempty" yaml:"defaultHeaders,omitempty"`
	MaxHeaderBytes      int               `json:"maxHeaderBytes,omitempty" yaml:"maxHeaderBytes,omitempty" env:"MAX_HEADER_BYTES"`
	ReadTimeout         time.Duration     `json:"readTimeout,omitempty" yaml:"readTimeout,omitempty" env:"READ_HEADER_TIMEOUT"`
	ReadHeaderTimeout   time.Duration     `json:"readHeaderTimeout,omitempty" yaml:"readHeaderTimeout,omitempty" env:"READ_HEADER_TIMEOUT"`
	WriteTimeout        time.Duration     `json:"writeTimeout,omitempty" yaml:"writeTimeout,omitempty" env:"WRITE_TIMEOUT"`
	IdleTimeout         time.Duration     `json:"idleTimeout,omitempty" yaml:"idleTimeout,omitempty" env:"IDLE_TIMEOUT"`
	ShutdownGracePeriod time.Duration     `json:"shutdownGracePeriod" yaml:"shutdownGracePeriod" env:"SHUTDOWN_GRACE_PERIOD"`

	Views ViewCacheConfig `json:"views,omitempty" yaml:"views,omitempty"`
}

// Resolve resolves the config from other sources.
func (c *Config) Resolve() error {
	return env.Env().ReadInto(c)
}

// BindAddrOrDefault returns the bind address or a default.
func (c Config) BindAddrOrDefault(defaults ...string) string {
	if len(c.BindAddr) > 0 {
		return c.BindAddr
	}
	if c.Port > 0 {
		return fmt.Sprintf(":%d", c.Port)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return DefaultBindAddr
}

// PortOrDefault returns the int32 port for a given config.
// This is useful in things like kubernetes pod templates.
// If the config .Port is unset, it will parse the .BindAddr,
// or the DefaultBindAddr for the port number.
func (c Config) PortOrDefault() int32 {
	if c.Port > 0 {
		return c.Port
	}
	if len(c.BindAddr) > 0 {
		return webutil.PortFromBindAddr(c.BindAddr)
	}
	return webutil.PortFromBindAddr(DefaultBindAddr)
}

// BaseURLOrDefault gets the base url for the app or a default.
func (c Config) BaseURLOrDefault() string {
	return c.BaseURL
}

// BaseURLIsSecureScheme returns if the base url starts with a secure scheme.
func (c Config) BaseURLIsSecureScheme() bool {
	if c.BaseURL == "" {
		return false
	}
	return strings.HasPrefix(strings.ToLower(c.BaseURL), SchemeHTTPS) || strings.HasPrefix(strings.ToLower(c.BaseURL), SchemeSPDY)
}

// SessionTimeoutOrDefault returns a property or a default.
func (c Config) SessionTimeoutOrDefault() time.Duration {
	if c.SessionTimeout > 0 {
		return c.SessionTimeout
	}
	return DefaultSessionTimeout
}

// CookieNameOrDefault returns a property or a default.
func (c Config) CookieNameOrDefault() string {
	if c.CookieName != "" {
		return c.CookieName
	}
	return DefaultCookieName
}

// CookiePathOrDefault returns a property or a default.
func (c Config) CookiePathOrDefault() string {
	if c.CookiePath != "" {
		return c.CookiePath
	}
	return DefaultCookiePath
}

// CookieDomainOrDefault returns a property or a default.
func (c Config) CookieDomainOrDefault() string {
	if c.CookieDomain != "" {
		return c.CookieDomain
	}
	return ""
}

// CookieSecureOrDefault returns a property or a default.
func (c Config) CookieSecureOrDefault() bool {
	if c.CookieSecure != nil {
		return *c.CookieSecure
	}
	if baseURL := c.BaseURLOrDefault(); baseURL != "" {
		return strings.HasPrefix(baseURL, SchemeHTTPS) || strings.HasPrefix(baseURL, SchemeSPDY)
	}
	return DefaultCookieSecure
}

// CookieHTTPOnlyOrDefault returns a property or a default.
func (c Config) CookieHTTPOnlyOrDefault() bool {
	if c.CookieHTTPOnly != nil {
		return *c.CookieHTTPOnly
	}
	return DefaultCookieHTTPOnly
}

// CookieSameSiteOrDefault returns a property or a default.
func (c Config) CookieSameSiteOrDefault() http.SameSite {
	if c.CookieSameSite != "" {
		return webutil.MustParseSameSite(c.CookieSameSite)
	}
	return 0
}

// MaxHeaderBytesOrDefault returns the maximum header size in bytes or a default.
func (c Config) MaxHeaderBytesOrDefault() int {
	if c.MaxHeaderBytes > 0 {
		return c.MaxHeaderBytes
	}
	return DefaultMaxHeaderBytes
}

// ReadTimeoutOrDefault gets a property.
func (c Config) ReadTimeoutOrDefault() time.Duration {
	if c.ReadTimeout > 0 {
		return c.ReadTimeout
	}
	return DefaultReadTimeout
}

// ReadHeaderTimeoutOrDefault gets a property.
func (c Config) ReadHeaderTimeoutOrDefault() time.Duration {
	if c.ReadHeaderTimeout > 0 {
		return c.ReadHeaderTimeout
	}
	return DefaultReadHeaderTimeout
}

// WriteTimeoutOrDefault gets a property.
func (c Config) WriteTimeoutOrDefault() time.Duration {
	if c.WriteTimeout > 0 {
		return c.WriteTimeout
	}
	return DefaultWriteTimeout
}

// IdleTimeoutOrDefault gets a property.
func (c Config) IdleTimeoutOrDefault() time.Duration {
	if c.IdleTimeout > 0 {
		return c.IdleTimeout
	}
	return DefaultIdleTimeout
}

// ShutdownGracePeriodOrDefault gets the shutdown grace period.
func (c Config) ShutdownGracePeriodOrDefault() time.Duration {
	if c.ShutdownGracePeriod > 0 {
		return c.ShutdownGracePeriod
	}
	return DefaultShutdownGracePeriod
}
