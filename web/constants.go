package web

import (
	"time"
)

const (
	// PackageName is the full name of this package.
	PackageName = "golang.blend.com/sdk/web"
	// RouteTokenFilepath is a special route token.
	RouteTokenFilepath = "filepath"
	// RegexpAssetCacheFiles is a common regex for parsing css, js, and html file routes.
	RegexpAssetCacheFiles = `^(.*)\.([0-9]+)\.(css|js|html|htm)$`
	// FieldTagPostForm is a field tag you can use to set a struct from a post body.
	FieldTagPostForm = "postForm"
)

const (
	// DefaultBindAddr is the default bind address.
	DefaultBindAddr = ":8080"
	// DefaultHealthzBindAddr is the default healthz bind address.
	DefaultHealthzBindAddr = ":8081"
	// DefaultMockBindAddr is a bind address used for integration testing.
	DefaultMockBindAddr = "127.0.0.1:0"
	// DefaultSkipRedirectTrailingSlash is the default if we should redirect for missing trailing slashes.
	DefaultSkipRedirectTrailingSlash = false
	// DefaultHandleOptions is a default.
	DefaultHandleOptions = false
	// DefaultHandleMethodNotAllowed is a default.
	DefaultHandleMethodNotAllowed = false
	// DefaultRecoverPanics returns if we should recover panics by default.
	DefaultRecoverPanics = true
	// DefaultMaxHeaderBytes is a default that is unset.
	DefaultMaxHeaderBytes = 0
	// DefaultReadTimeout is a default.
	DefaultReadTimeout = 0
	// DefaultReadHeaderTimeout is a default.
	DefaultReadHeaderTimeout time.Duration = 0
	// DefaultWriteTimeout is a default.
	DefaultWriteTimeout time.Duration = 0
	// DefaultIdleTimeout is a default.
	DefaultIdleTimeout time.Duration = 0
	// DefaultCookieName is the default name of the field that contains the session id.
	DefaultCookieName = "SID"
	// DefaultSecureCookieName is the default name of the field that contains the secure session id.
	DefaultSecureCookieName = "SSID"
	// DefaultCookiePath is the default cookie path.
	DefaultCookiePath = "/"
	// DefaultCookieSecure returns what the default value for the `Secure` bit of issued cookies will be.
	DefaultCookieSecure = true
	// DefaultCookieHTTPOnly returns what the default value for the `HTTPOnly` bit of issued cookies will be.
	DefaultCookieHTTPOnly = true
	// DefaultSessionTimeout is the default absolute timeout for a session (24 hours as a sane default).
	DefaultSessionTimeout time.Duration = 24 * time.Hour
	// DefaultUseSessionCache is the default if we should use the auth manager session cache.
	DefaultUseSessionCache = true
	// DefaultSessionTimeoutIsAbsolute is the default if we should set absolute session expiries.
	DefaultSessionTimeoutIsAbsolute = true
	// DefaultHTTPSUpgradeTargetPort is the default upgrade target port.
	DefaultHTTPSUpgradeTargetPort = 443
	// DefaultKeepAlive is the default setting for TCP KeepAlive.
	DefaultKeepAlive = true
	// DefaultKeepAlivePeriod is the default time to keep tcp connections open.
	DefaultKeepAlivePeriod = 3 * time.Minute
	// DefaultShutdownGracePeriod is the default shutdown grace period.
	DefaultShutdownGracePeriod = 30 * time.Second
	// DefaultHealthzFailureThreshold is the default healthz failure threshold.
	DefaultHealthzFailureThreshold = 3
	// DefaultViewBufferPoolSize is the default buffer pool size.
	DefaultViewBufferPoolSize = 256
)

const (
	// LenSessionID is the byte length of a session id.
	LenSessionID = 64
	// LenSessionIDBase64 is the length of a session id base64 encoded.
	LenSessionIDBase64 = 88
)
