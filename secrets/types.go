package secrets

import "time"

// Values is a bag of values.
type Values = map[string]string

// SecretV1 is the structure returned for every secret within Vault.
type SecretV1 struct {
	// The request ID that generated this response
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`
	// Data is the actual contents of the secret. The format of the data
	// is arbitrary and up to the secret backend.
	Data Values `json:"data"`
	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings []string `json:"warnings"`
	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth *SecretAuth `json:"auth,omitempty"`
	// WrapInfo, if non-nil, means that the initial response was wrapped in the
	// cubbyhole of the given token (which has a TTL of the given number of
	// seconds)
	WrapInfo *SecretWrapInfo `json:"wrap_info,omitempty"`
}

// SecretListV1 is the structure returned for a list of secret keys in vault
type SecretListV1 struct {
	// The request ID that generated this response
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`
	// Data is the list of keys and subfolders at this path. Subfolders end with a slash, keys do not
	Data KeyData `json:"data"`
	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings []string `json:"warnings"`
	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth *SecretAuth `json:"auth,omitempty"`
	// WrapInfo, if non-nil, means that the initial response was wrapped in the
	// cubbyhole of the given token (which has a TTL of the given number of
	// seconds)
	WrapInfo *SecretWrapInfo `json:"wrap_info,omitempty"`
}

// SecretV2 is the structure returned for every secret within Vault.
type SecretV2 struct {
	// The request ID that generated this response
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`
	// Data is the actual contents of the secret. The format of the data
	// is arbitrary and up to the secret backend.
	Data SecretData `json:"data"`
	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings []string `json:"warnings"`
	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth *SecretAuth `json:"auth,omitempty"`
	// WrapInfo, if non-nil, means that the initial response was wrapped in the
	// cubbyhole of the given token (which has a TTL of the given number of
	// seconds)
	WrapInfo *SecretWrapInfo `json:"wrap_info,omitempty"`
}

// SecretListV2 is the structure returned for every secret within Vault.
type SecretListV2 struct {
	// The request ID that generated this response
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`
	// Data is the list of keys and subfolders at this path. Subfolders end with a slash, keys do not
	Data KeyData `json:"data"`
	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings []string `json:"warnings"`
	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth *SecretAuth `json:"auth,omitempty"`
	// WrapInfo, if non-nil, means that the initial response was wrapped in the
	// cubbyhole of the given token (which has a TTL of the given number of
	// seconds)
	WrapInfo *SecretWrapInfo `json:"wrap_info,omitempty"`
}

// TransitKey is the structure returned for every transit key within Vault.
type TransitKey struct {
	// The request ID that generated this response
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`
	// Data is the data associated with a transit key
	Data map[string]interface{} `json:"data"`
	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings []string `json:"warnings"`
	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth *SecretAuth `json:"auth,omitempty"`
	// WrapInfo, if non-nil, means that the initial response was wrapped in the
	// cubbyhole of the given token (which has a TTL of the given number of
	// seconds)
	WrapInfo *SecretWrapInfo `json:"wrap_info,omitempty"`
}

// SecretData is used for puts.
type SecretData struct {
	Data Values `json:"data"`
}

// KeyData is used for lists.
type KeyData struct {
	Keys []string `json:"keys"`
}

// SecretAuth is the structure containing auth information if we have it.
type SecretAuth struct {
	ClientToken   string            `json:"client_token"`
	Accessor      string            `json:"accessor"`
	Policies      []string          `json:"policies"`
	Metadata      map[string]string `json:"metadata"`
	LeaseDuration int               `json:"lease_duration"`
	Renewable     bool              `json:"renewable"`
}

// SecretWrapInfo contains wrapping information if we have it. If what is
// contained is an authentication token, the accessor for the token will be
// available in WrappedAccessor.
type SecretWrapInfo struct {
	Token           string    `json:"token"`
	Accessor        string    `json:"accessor"`
	TTL             int       `json:"ttl"`
	CreationTime    time.Time `json:"creation_time"`
	CreationPath    string    `json:"creation_path"`
	WrappedAccessor string    `json:"wrapped_accessor"`
}

// MountResponse is the result of a call to a mount.
type MountResponse struct {
	RequestID string `json:"request_id"`
	Data      Mount  `json:"data"`
}

// Mount is a vault mount.
type Mount struct {
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Accessor    string            `json:"accessor"`
	Config      MountConfig       `json:"config"`
	Options     map[string]string `json:"options"`
	Local       bool              `json:"local"`
	SealWrap    bool              `json:"seal_wrap" mapstructure:"seal_wrap"`
}

// MountConfig is a vault mount config.
type MountConfig struct {
	DefaultLeaseTTL           int      `json:"default_lease_ttl" mapstructure:"default_lease_ttl"`
	MaxLeaseTTL               int      `json:"max_lease_ttl" mapstructure:"max_lease_ttl"`
	ForceNoCache              bool     `json:"force_no_cache" mapstructure:"force_no_cache"`
	PluginName                string   `json:"plugin_name,omitempty" mapstructure:"plugin_name"`
	AuditNonHMACRequestKeys   []string `json:"audit_non_hmac_request_keys,omitempty" mapstructure:"audit_non_hmac_request_keys"`
	AuditNonHMACResponseKeys  []string `json:"audit_non_hmac_response_keys,omitempty" mapstructure:"audit_non_hmac_response_keys"`
	ListingVisibility         string   `json:"listing_visibility,omitempty" mapstructure:"listing_visibility"`
	PassthroughRequestHeaders []string `json:"passthrough_request_headers,omitempty" mapstructure:"passthrough_request_headers"`
}

// MountInput is a vault mount input.
type MountInput struct {
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Config      MountConfigInput  `json:"config"`
	Options     map[string]string `json:"options"`
	Local       bool              `json:"local"`
	PluginName  string            `json:"plugin_name,omitempty"`
	SealWrap    bool              `json:"seal_wrap" mapstructure:"seal_wrap"`
}

// MountConfigInput is a vault mount config input.
type MountConfigInput struct {
	Options                   map[string]string `json:"options" mapstructure:"options"`
	DefaultLeaseTTL           string            `json:"default_lease_ttl" mapstructure:"default_lease_ttl"`
	MaxLeaseTTL               string            `json:"max_lease_ttl" mapstructure:"max_lease_ttl"`
	ForceNoCache              bool              `json:"force_no_cache" mapstructure:"force_no_cache"`
	PluginName                string            `json:"plugin_name,omitempty" mapstructure:"plugin_name"`
	AuditNonHMACRequestKeys   []string          `json:"audit_non_hmac_request_keys,omitempty" mapstructure:"audit_non_hmac_request_keys"`
	AuditNonHMACResponseKeys  []string          `json:"audit_non_hmac_response_keys,omitempty" mapstructure:"audit_non_hmac_response_keys"`
	ListingVisibility         string            `json:"listing_visibility,omitempty" mapstructure:"listing_visibility"`
	PassthroughRequestHeaders []string          `json:"passthrough_request_headers,omitempty" mapstructure:"passthrough_request_headers"`
}

// TransitResult is the structure returned by vault for transit requests
type TransitResult struct {
	Data struct {
		Ciphertext string `json:"ciphertext"`
		Plaintext  string `json:"plaintext"`
	} `json:"data"`
}
