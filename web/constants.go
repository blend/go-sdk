package web

import "time"

const (
	// PackageName is the full name of this package.
	PackageName = "github.com/blend/go-sdk/web"

	// HeaderAllow is a common header.
	HeaderAllow = "Allow"

	// RouteTokenFilepath is a special route token.
	RouteTokenFilepath = "filepath"

	// RegexpAssetCacheFiles is a common regex for parsing css, js, and html file routes.
	RegexpAssetCacheFiles = `^(.*)\.([0-9]+)\.(css|js|html|htm)$`

	// HeaderAcceptEncoding is the "Accept-Encoding" header.
	// It indicates what types of encodings the request will accept responses as.
	// It typically enables or disables compressed (gzipped) responses.
	HeaderAcceptEncoding = "Accept-Encoding"

	// HeaderSetCookie is the header that sets cookies in a response.
	HeaderSetCookie = "Set-Cookie"

	// HeaderCookie is the request cookie header.
	HeaderCookie = "Cookie"

	// HeaderDate is the "Date" header.
	// It provides a timestamp the response was generated at.
	// It is typically used by client cache control to invalidate expired items.
	HeaderDate = "Date"

	// HeaderCacheControl is the "Cache-Control" header.
	// It indicates if and how clients should cache responses.
	// Typical values for this include "no-cache", "max-age", "min-fresh", and "max-stale" variants.
	HeaderCacheControl = "Cache-Control"

	// HeaderConnection is the "Connection" header.
	// It is used to indicate if the connection should remain open by the server
	// after the final response bytes are sent.
	// This allows the connection to be re-used, helping mitigate connection negotiation
	// penalites in making requests.
	HeaderConnection = "Connection"

	// HeaderContentEncoding is the "Content-Encoding" header.
	// It is used to indicate what the response encoding is.
	// Typical values are "gzip", "deflate", "compress", "br", and "identity" indicating no compression.
	HeaderContentEncoding = "Content-Encoding"

	// HeaderContentLength is the "Content-Length" header.
	// If provided, it specifies the size of the request or response.
	HeaderContentLength = "Content-Length"

	// HeaderContentType is the "Content-Type" header.
	// It specifies the MIME-type of the request or response.
	HeaderContentType = "Content-Type"

	// HeaderServer is the "Server" header.
	// It is an informational header to tell the client what server software was used.
	HeaderServer = "Server"

	// HeaderUserAgent is the user agent header.
	HeaderUserAgent = "User-Agent"

	// HeaderETag is the etag header.
	HeaderETag = "etag"

	// HeaderVary is the "Vary" header.
	// It is used to indicate what fields should be used by the client as cache keys.
	HeaderVary = "Vary"

	// HeaderXServedBy is the "X-Served-By" header.
	// It is an informational header that indicates what software was used to generate the response.
	HeaderXServedBy = "X-Served-By"

	// HeaderXFrameOptions is the "X-Frame-Options" header.
	// It indicates if a browser is allowed to render the response in a <frame> element or not.
	HeaderXFrameOptions = "X-Frame-Options"

	// HeaderXXSSProtection is the "X-Xss-Protection" header.
	// It is a feature of internet explorer, and indicates if the browser should allow
	// requests across domain boundaries.
	HeaderXXSSProtection = "X-Xss-Protection"

	// HeaderXContentTypeOptions is the "X-Content-Type-Options" header.
	HeaderXContentTypeOptions = "X-Content-Type-Options"

	// HeaderStrictTransportSecurity is the hsts header.
	HeaderStrictTransportSecurity = "Strict-Transport-Security"

	// ContentTypeApplicationJSON is a content type for JSON responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeApplicationJSON = "application/json; charset=UTF-8"

	// ContentTypeHTML is a content type for html responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeHTML = "text/html; charset=utf-8"

	//ContentTypeXML is a content type for XML responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeXML = "text/xml; charset=utf-8"

	// ContentTypeText is a content type for text responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeText = "text/plain; charset=utf-8"

	// ConnectionKeepAlive is a value for the "Connection" header and
	// indicates the server should keep the tcp connection open
	// after the last byte of the response is sent.
	ConnectionKeepAlive = "keep-alive"

	// ContentEncodingIdentity is the identity (uncompressed) content encoding.
	ContentEncodingIdentity = "identity"
	// ContentEncodingGZIP is the gzip (compressed) content encoding.
	ContentEncodingGZIP = "gzip"
)

const (
	// SchemeHTTP is a protocol scheme.
	SchemeHTTP = "http"

	// SchemeHTTPS is a protocol scheme.
	SchemeHTTPS = "https"

	// SchemeSPDY is a protocol scheme.
	SchemeSPDY = "spdy"
)

const (
	// MethodGet is an http verb.
	MethodGet = "GET"

	// MethodPost is an http verb.
	MethodPost = "POST"

	// MethodPut is an http verb.
	MethodPut = "PUT"

	// MethodDelete is an http verb.
	MethodDelete = "DELETE"

	// MethodConnect is an http verb.
	MethodConnect = "CONNECT"

	// MethodOptions is an http verb.
	MethodOptions = "OPTIONS"
)

const (
	// HSTSMaxAgeFormat is the format string for a max age token.
	HSTSMaxAgeFormat = "max-age=%d"

	// HSTSIncludeSubDomains is a header value token.
	HSTSIncludeSubDomains = "includeSubDomains"

	// HSTSPreload is a header value token.
	HSTSPreload = "preload"
)

// Environment Variables
const (
	// EnvironmentVariableBindAddr is an env var that determines (if set) what the bind address should be.
	EnvironmentVariableBindAddr = "BIND_ADDR"

	// EnvironmentVariableHealthzBindAddr is an env var that determines (if set) what the healthz sidecar bind address should be.
	EnvironmentVariableHealthzBindAddr = "HEALTHZ_BIND_ADDR"

	// EnvironmentVariableUpgraderBindAddr is an env var that determines (if set) what the bind address should be.
	EnvironmentVariableUpgraderBindAddr = "UPGRADER_BIND_ADDR"

	// EnvironmentVariablePort is an env var that determines what the default bind address port segment returns.
	EnvironmentVariablePort = "PORT"

	// EnvironmentVariableHealthzPort is an env var that determines what the default healthz bind address port segment returns.
	EnvironmentVariableHealthzPort = "HEALTHZ_PORT"

	// EnvironmentVariableUpgraderPort is an env var that determines what the default bind address port segment returns.
	EnvironmentVariableUpgraderPort = "UPGRADER_PORT"

	// EnvironmentVariableTLSCert is an env var that contains the TLS cert.
	EnvironmentVariableTLSCert = "TLS_CERT"

	// EnvironmentVariableTLSKey is an env var that contains the TLS key.
	EnvironmentVariableTLSKey = "TLS_KEY"

	// EnvironmentVariableTLSCertFile is an env var that contains the file path to the TLS cert.
	EnvironmentVariableTLSCertFile = "TLS_CERT_FILE"

	// EnvironmentVariableTLSKeyFile is an env var that contains the file path to the TLS key.
	EnvironmentVariableTLSKeyFile = "TLS_KEY_FILE"
)

// Defaults
const (
	// DefaultBindAddr is the default bind address.
	DefaultBindAddr = ":8080"
	// DefaultHealthzBindAddr is the default healthz bind address.
	DefaultHealthzBindAddr = ":8081"
	// DefaultIntegrationBindAddr is a bind address used for integration testing.
	DefaultIntegrationBindAddr = "127.0.0.1:0"
	// DefaultRedirectTrailingSlash is the default if we should redirect for missing trailing slashes.
	DefaultRedirectTrailingSlash = true
	// DefaultHandleOptions is a default.
	DefaultHandleOptions = false
	// DefaultHandleMethodNotAllowed is a default.
	DefaultHandleMethodNotAllowed = false
	// DefaultRecoverPanics returns if we should recover panics by default.
	DefaultRecoverPanics = true

	// DefaultHSTS is the default for if hsts is enabled.
	DefaultHSTS = true
	// DefaultHSTSMaxAgeSeconds is the default hsts max age seconds.
	DefaultHSTSMaxAgeSeconds = 31536000
	// DefaultHSTSIncludeSubDomains is a default.
	DefaultHSTSIncludeSubDomains = true
	// DefaultHSTSPreload is a default.
	DefaultHSTSPreload = true
	// DefaultMaxHeaderBytes is a default that is unset.
	DefaultMaxHeaderBytes = 0
	// DefaultReadTimeout is a default.
	DefaultReadTimeout = 5 * time.Second
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
	// DefaultSessionTimeout is the default absolute timeout for a session (24 hours as a sane default).
	DefaultSessionTimeout time.Duration = 24 * time.Hour
	// DefaultUseSessionCache is the default if we should use the auth manager session cache.
	DefaultUseSessionCache = true
	// DefaultSessionTimeoutIsAbsolute is the default if we should set absolute session expiries.
	DefaultSessionTimeoutIsAbsolute = true

	// DefaultHTTPSUpgradeTargetPort is the default upgrade target port.
	DefaultHTTPSUpgradeTargetPort = 443

	// DefaultShutdownGracePeriod is the default shutdown grace period.
	DefaultShutdownGracePeriod = 30 * time.Second

	// DefaultHealthzFailureThreshold is the default healthz failure threshold.
	DefaultHealthzFailureThreshold = 3

	// DefaultBufferPoolSize is the default buffer pool size.
	DefaultViewBufferPoolSize = 256
)

// DefaultHeaders are the default headers added by go-web.
var DefaultHeaders = map[string]string{
	HeaderServer:    PackageName,
	HeaderXServedBy: PackageName,
}

// SessionLockPolicy is a lock policy.
type SessionLockPolicy int

const (
	// SessionUnsafe is a lock-free session policy.
	SessionUnsafe SessionLockPolicy = 0

	// SessionReadLock is a lock policy that acquires a read lock on session.
	SessionReadLock SessionLockPolicy = 1

	// SessionReadWriteLock is a lock policy that acquires both a read and a write lock on session.
	SessionReadWriteLock SessionLockPolicy = 2
)

const (
	// PostBodySize is the maximum post body size we will typically consume.
	PostBodySize = int64(1 << 26) //64mb

	// PostBodySizeMax is the absolute maximum file size the server can handle.
	PostBodySizeMax = int64(1 << 32) //enormous.
)

const (
	// LenSessionID is the byte length of a session id.
	LenSessionID = 64
	// LenSessionIDBase64 is the length of a session id base64 encoded.
	LenSessionIDBase64 = 88
)

// test keys
const (
	TestTLSCert = `-----BEGIN CERTIFICATE-----
MIIC+jCCAeKgAwIBAgIRAKGQgEUjhTZMM2VMx9y92MUwDQYJKoZIhvcNAQELBQAw
EjEQMA4GA1UEChMHQWNtZSBDbzAeFw0xODAzMDkwMzIxMzhaFw0xOTAzMDkwMzIx
MzhaMBIxEDAOBgNVBAoTB0FjbWUgQ28wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDCONjExGZ+MwYZ1CosUB+sa9jS/AD0YkOi8AgiOYughLrKx5RuSsO9
ZaO0iwH987SFwAxBEiXwfLceEDgHYLGNfKQdYMCdh1yclr9yKrfpLV1SvPwT/utm
ek3ONwbJwqIrBP0dNWtfRhHhu2Gyc1JjxpqETdCUUZfuJWouVjVIxaIxLvyxYkUo
AS6SpUlUOOF3Wnre4+3x1RWRpXwns/HUFjsQBOIBo7pganxzcukTsQZWv+kJEA2o
EW33VdLQBuD59X6h1/qjx93s3AndeT5CoeVCAQ6PKXuV9z1WCpRewPpD+J89Noff
aueXIhTvxpFnB6W6VGVDQmnhEbnwA2IPAgMBAAGjSzBJMA4GA1UdDwEB/wQEAwIF
oDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMBQGA1UdEQQNMAuC
CWxvY2FsaG9zdDANBgkqhkiG9w0BAQsFAAOCAQEAYkkoNdditdKaEWrUjMc52QqJ
e4hbjqWT6W3bphGgYiKvnxgcDQYL3+RgEd7tGIHfgLkIiuM9efH+KJ4/jdXFWlcQ
7PoS9nGn0FwNvGdt9KCzNZSODSgQNt7FdsSpfw6Qzhn6XCwx3Bay9uF6cPap+wtX
SX6fD+az+dh0UPYoEltuKBv43+wLwsxAg18vBFuACI52NomvNw4G4uw4epBGGmp8
A0A4h9O67T/bFXchS+uIQnThZo4U/TCDu0xi/Q89xtjWff1YybwR85l85pEt1v7G
ei1eKWKYUxUU7lBMaECknLsJ4xsDKRSA5tvEDCkeQDCwTD7Msh5uGQ9itoWMlQ==
-----END CERTIFICATE-----
`
	TestTLSKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwjjYxMRmfjMGGdQqLFAfrGvY0vwA9GJDovAIIjmLoIS6yseU
bkrDvWWjtIsB/fO0hcAMQRIl8Hy3HhA4B2CxjXykHWDAnYdcnJa/ciq36S1dUrz8
E/7rZnpNzjcGycKiKwT9HTVrX0YR4bthsnNSY8aahE3QlFGX7iVqLlY1SMWiMS78
sWJFKAEukqVJVDjhd1p63uPt8dUVkaV8J7Px1BY7EATiAaO6YGp8c3LpE7EGVr/p
CRANqBFt91XS0Abg+fV+odf6o8fd7NwJ3Xk+QqHlQgEOjyl7lfc9VgqUXsD6Q/if
PTaH32rnlyIU78aRZwelulRlQ0Jp4RG58ANiDwIDAQABAoIBAHcg8yTN6qfhmA5j
qnJ/us3BYL8Yv2UmmKHqZLLJZTFR+FjEzfBQf3s+SolE8jXYM5QOVfXbsdWuSYtx
G0y7LGzCVM+INtzo2A9cD5VxSlkF8EX9kQiaxbyXq/2eltVOQrXsW2x9BZzsl69D
hgs03QZCHSilqhgva+cwn85IJmq5bL5BMlNT1vFUgKz4QWISuBQc84PpH9R3P0oF
ur4PRJuh6Q3/GX2MF7fuNw+cweg6lNM2IlVmoH3jJo4byW+tzruv5O+/0s92CsSM
s5ywkZlgydrh1w4Irqli67y/jdDdA9zHcr+DBpVquJ1arez/ImRtKA9+FRNP4YvM
k3FOh0ECgYEA2UU+8iad7Kd7bcrhCq6AItlv51MxTp9ASDoFiCFJncTOGLdzcVNA
a+reF22XYdD32R94ldWGlIBp3MbNTyK5HYkTbwHG8414fahxg3Uy4je0NLQzHpIH
OQjaX+YFUtMDaGL7MCIDeC1FKCwfnWBRS/6xaZe3g4ne1wqZZ46DmxECgYEA5NfT
jsLSPXD5ZEz594jsOfTJ24RH4CgB69BQTd9z9AezMlTZE3fTUeXjhZRim1cs/+/4
lotnMuUEYOVRwtfJS+hqVGg1y7MFJMTo7O5RP2+SynnIrXBkZgXwKfNX0Dj4crnA
dxlUHPPFzNEZzkNMDuiwo4ERs5G+11OPD+UL6x8CgYAkaUVmQXB/44V83d4e8yWI
MZZeVwPRYEDemdKpgKKcrQm4/K19FW2baE318SjIfMO8gFiuC421P1v+YtavZ2tM
dtdp6AtWb6P8swjq9e4kGR+7IWPbwK8zMLegEKVdvv04NjZQV7LrJfMMC3D059pX
+QP0ZTec9LMCqMUSpMCLcQKBgQDGnjAnGx6AZzp9fHYECxoEX1qHpTMA8ZhhRGc+
f2/TYI9+YrgZtol57o5f1N8Utj//TxcyCoIiYTVAqCgjdUhoEque4Oe4CYOwWxtS
8LEh3sPH6pVrOz5YclT1BBi2R4wTfvb2J8yiaE3IK8A7DpvH4NvWvWJQuXGq0AI+
KG0EvwKBgB8nHRWRbNJ8admJukGb5HF2mS1tDuHi+vB1dsTydfPDyf33B1HoEG0p
mfr9uzS9ndAYCopZO33b1h65wlPP6jnIJheycn15n7HRjYezTr8cODMnJLrRotAJ
HCsYkCmGXiwJN2guZo6l/5+GqRo3SN19dZptrH/rC/wAai0+Ctqw
-----END RSA PRIVATE KEY-----`
)
