/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package vault

import "time"

// Values is a bag of values.
type Values = map[string]interface{}

// SecretV1 is the structure returned for every secret within Vault.
type SecretV1 struct {
	// The request ID that generated this response
	RequestID	string	`json:"request_id"`
	LeaseID		string	`json:"lease_id"`
	LeaseDuration	int	`json:"lease_duration"`
	Renewable	bool	`json:"renewable"`
	// Data is the actual contents of the secret. The format of the data
	// is arbitrary and up to the secret backend.
	Data	Values	`json:"data"`
	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings	[]string	`json:"warnings"`
	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth	*SecretAuth	`json:"auth,omitempty"`
	// WrapInfo, if non-nil, means that the initial response was wrapped in the
	// cubbyhole of the given token (which has a TTL of the given number of
	// seconds)
	WrapInfo	*SecretWrapInfo	`json:"wrap_info,omitempty"`
}

// SecretListV1 is the structure returned for a list of secret keys in vault
type SecretListV1 struct {
	// The request ID that generated this response
	RequestID	string	`json:"request_id"`
	LeaseID		string	`json:"lease_id"`
	LeaseDuration	int	`json:"lease_duration"`
	Renewable	bool	`json:"renewable"`
	// Data is the list of keys and subfolders at this path. Subfolders end with a slash, keys do not
	Data	KeyData	`json:"data"`
	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings	[]string	`json:"warnings"`
	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth	*SecretAuth	`json:"auth,omitempty"`
	// WrapInfo, if non-nil, means that the initial response was wrapped in the
	// cubbyhole of the given token (which has a TTL of the given number of
	// seconds)
	WrapInfo	*SecretWrapInfo	`json:"wrap_info,omitempty"`
}

// SecretV2 is the structure returned for every secret within Vault.
type SecretV2 struct {
	// The request ID that generated this response
	RequestID	string	`json:"request_id"`
	LeaseID		string	`json:"lease_id"`
	LeaseDuration	int	`json:"lease_duration"`
	Renewable	bool	`json:"renewable"`
	// Data is the actual contents of the secret. The format of the data
	// is arbitrary and up to the secret backend.
	Data	SecretData	`json:"data"`
	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings	[]string	`json:"warnings"`
	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth	*SecretAuth	`json:"auth,omitempty"`
	// WrapInfo, if non-nil, means that the initial response was wrapped in the
	// cubbyhole of the given token (which has a TTL of the given number of
	// seconds)
	WrapInfo	*SecretWrapInfo	`json:"wrap_info,omitempty"`
}

// SecretListV2 is the structure returned for every secret within Vault.
type SecretListV2 struct {
	// The request ID that generated this response
	RequestID	string	`json:"request_id"`
	LeaseID		string	`json:"lease_id"`
	LeaseDuration	int	`json:"lease_duration"`
	Renewable	bool	`json:"renewable"`
	// Data is the list of keys and subfolders at this path. Subfolders end with a slash, keys do not
	Data	KeyData	`json:"data"`
	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings	[]string	`json:"warnings"`
	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth	*SecretAuth	`json:"auth,omitempty"`
	// WrapInfo, if non-nil, means that the initial response was wrapped in the
	// cubbyhole of the given token (which has a TTL of the given number of
	// seconds)
	WrapInfo	*SecretWrapInfo	`json:"wrap_info,omitempty"`
}

// TransitKey is the structure returned for every transit key within Vault.
type TransitKey struct {
	// The request ID that generated this response
	RequestID	string	`json:"request_id"`
	LeaseID		string	`json:"lease_id"`
	LeaseDuration	int	`json:"lease_duration"`
	Renewable	bool	`json:"renewable"`
	// Data is the data associated with a transit key
	Data	map[string]interface{}	`json:"data"`
	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings	[]string	`json:"warnings"`
	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth	*SecretAuth	`json:"auth,omitempty"`
	// WrapInfo, if non-nil, means that the initial response was wrapped in the
	// cubbyhole of the given token (which has a TTL of the given number of
	// seconds)
	WrapInfo	*SecretWrapInfo	`json:"wrap_info,omitempty"`
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
	ClientToken	string			`json:"client_token"`
	Accessor	string			`json:"accessor"`
	Policies	[]string		`json:"policies"`
	Metadata	map[string]string	`json:"metadata"`
	LeaseDuration	int			`json:"lease_duration"`
	Renewable	bool			`json:"renewable"`
}

// SecretWrapInfo contains wrapping information if we have it. If what is
// contained is an authentication token, the accessor for the token will be
// available in WrappedAccessor.
type SecretWrapInfo struct {
	Token		string		`json:"token"`
	Accessor	string		`json:"accessor"`
	TTL		int		`json:"ttl"`
	CreationTime	time.Time	`json:"creation_time"`
	CreationPath	string		`json:"creation_path"`
	WrappedAccessor	string		`json:"wrapped_accessor"`
}

// MountResponse is the result of a call to a mount.
type MountResponse struct {
	RequestID	string	`json:"request_id"`
	Data		Mount	`json:"data"`
}

// Mount is a vault mount.
type Mount struct {
	Type		string			`json:"type"`
	Description	string			`json:"description"`
	Accessor	string			`json:"accessor"`
	Config		MountConfig		`json:"config"`
	Options		map[string]string	`json:"options"`
	Local		bool			`json:"local"`
	SealWrap	bool			`json:"seal_wrap" mapstructure:"seal_wrap"`
}

// MountConfig is a vault mount config.
type MountConfig struct {
	DefaultLeaseTTL			int		`json:"default_lease_ttl" mapstructure:"default_lease_ttl"`
	MaxLeaseTTL			int		`json:"max_lease_ttl" mapstructure:"max_lease_ttl"`
	ForceNoCache			bool		`json:"force_no_cache" mapstructure:"force_no_cache"`
	PluginName			string		`json:"plugin_name,omitempty" mapstructure:"plugin_name"`
	AuditNonHMACRequestKeys		[]string	`json:"audit_non_hmac_request_keys,omitempty" mapstructure:"audit_non_hmac_request_keys"`
	AuditNonHMACResponseKeys	[]string	`json:"audit_non_hmac_response_keys,omitempty" mapstructure:"audit_non_hmac_response_keys"`
	ListingVisibility		string		`json:"listing_visibility,omitempty" mapstructure:"listing_visibility"`
	PassthroughRequestHeaders	[]string	`json:"passthrough_request_headers,omitempty" mapstructure:"passthrough_request_headers"`
}

// MountInput is a vault mount input.
type MountInput struct {
	Type		string			`json:"type"`
	Description	string			`json:"description"`
	Config		MountConfigInput	`json:"config"`
	Options		map[string]string	`json:"options"`
	Local		bool			`json:"local"`
	PluginName	string			`json:"plugin_name,omitempty"`
	SealWrap	bool			`json:"seal_wrap" mapstructure:"seal_wrap"`
}

// MountConfigInput is a vault mount config input.
type MountConfigInput struct {
	Options				map[string]string	`json:"options" mapstructure:"options"`
	DefaultLeaseTTL			string			`json:"default_lease_ttl" mapstructure:"default_lease_ttl"`
	MaxLeaseTTL			string			`json:"max_lease_ttl" mapstructure:"max_lease_ttl"`
	ForceNoCache			bool			`json:"force_no_cache" mapstructure:"force_no_cache"`
	PluginName			string			`json:"plugin_name,omitempty" mapstructure:"plugin_name"`
	AuditNonHMACRequestKeys		[]string		`json:"audit_non_hmac_request_keys,omitempty" mapstructure:"audit_non_hmac_request_keys"`
	AuditNonHMACResponseKeys	[]string		`json:"audit_non_hmac_response_keys,omitempty" mapstructure:"audit_non_hmac_response_keys"`
	ListingVisibility		string			`json:"listing_visibility,omitempty" mapstructure:"listing_visibility"`
	PassthroughRequestHeaders	[]string		`json:"passthrough_request_headers,omitempty" mapstructure:"passthrough_request_headers"`
}

// BatchTransitInput is the structure of batch encrypt / decrypt requests
type BatchTransitInput struct {
	BatchTransitInputItems []BatchTransitInputItem `json:"batch_input"`
}

// BatchTransitInputItem is a single item in a batch encrypt / decrypt request
type BatchTransitInputItem struct {
	Context		[]byte	`json:"context,omitempty"`
	Ciphertext	string	`json:"ciphertext,omitempty"`
	Plaintext	[]byte	`json:"plaintext,omitempty"`
}

// BatchTransitResult is the structure returned by vault for batch transit requests
type BatchTransitResult struct {
	Data struct {
		BatchTransitResult []struct {
			// Error, if set represents a failure encountered while encrypting/decrypting a
			// corresponding batch request item
			Error		string	`json:"error"`
			Ciphertext	string	`json:"ciphertext"`
			Plaintext	string	`json:"plaintext"`
		} `json:"batch_results"`
	} `json:"data"`
}

// TransitResult is the structure returned by vault for transit requests
type TransitResult struct {
	Data struct {
		Ciphertext	string	`json:"ciphertext"`
		Plaintext	string	`json:"plaintext"`
	} `json:"data"`
}

// TransitHmacResult is the structure returned by vault for transit hmac requests
type TransitHmacResult struct {
	Data struct {
		Hmac string `json:"hmac"`
	} `json:"data"`
}

// CreateTransitKeyConfig is the configuration data for creating a TransitKey
type CreateTransitKeyConfig struct {
	// Convergent - If enabled, the key will support convergent encryption, where the same plaintext creates the same
	// ciphertext. This requires derived to be set to true. When enabled, each encryption(/decryption/rewrap/datakey)
	// operation will derive a nonce value rather than randomly generate it.
	Convergent	bool	`json:"convergent_encryption,omitempty"`
	// Derived - Specifies if key derivation is to be used. If enabled, all encrypt/decrypt requests to this named key
	// must provide a context which is used for key derivation.
	Derived	bool	`json:"derived,omitempty"`
	// Exportable - Enables keys to be exportable. This allows for all the valid keys in the key ring to be exported.
	// Once set, this cannot be disabled.
	Exportable	bool	`json:"exportable,omitempty"`
	// AllowPlaintextBackup - If set, enables taking backup of named key in the plaintext format. Once set, this cannot
	// be disabled.
	AllowPlaintextBackup	bool	`json:"allow_plaintext_backup,omitempty"`
	// Type specifies the type of key to create. The default type is "aes256-gcm96":
	//   aes256-gcm96 – AES-256 wrapped with GCM using a 96-bit nonce size AEAD (symmetric, supports derivation and
	//      convergent encryption)
	//   chacha20-poly1305 – ChaCha20-Poly1305 AEAD (symmetric, supports derivation and convergent encryption)
	//   ed25519 – ED25519 (asymmetric, supports derivation). When using derivation, a sign operation with the same
	//      context will derive the same key and signature; this is a signing analog to convergent_encryption.
	//   ecdsa-p256 – ECDSA using the P-256 elliptic curve (asymmetric)
	//   rsa-2048 - RSA with bit size of 2048 (asymmetric)
	//   rsa-4096 - RSA with bit size of 4096 (asymmetric)
	Type	string	`json:"type,omitempty"`
}

// UpdateTransitKeyConfig is the configuration data for modifying a TransitKey
type UpdateTransitKeyConfig struct {
	// MinDecryptionVersion -  Specifies the minimum version of ciphertext allowed to be decrypted. Adjusting this as
	// part of a key rotation policy can prevent old copies of ciphertext from being decrypted, should they fall into
	// the wrong hands. For signatures, this value controls the minimum version of signature that can be verified
	// against. For HMACs, this controls the minimum version of a key allowed to be used as the key for verification.
	MinDecryptionVersion	int	`json:"min_decryption_version,omitempty"`
	// MinEncryptionVersion - Specifies the minimum version of the key that can be used to encrypt plaintext, sign
	// payloads, or generate HMACs. Must be 0 (which will use the latest version) or a value greater or equal to
	// min_decryption_version.
	MinEncryptionVersion	int	`json:"min_encryption_version,omitempty"`
	// DeletionAllowed - Specifies if the key is allowed to be deleted.
	DeletionAllowed	*bool	`json:"deletion_allowed,omitempty"`
	// Exportable - Enables keys to be exportable. This allows for all the valid keys in the key ring to be exported.
	// Once set, this cannot be disabled.
	Exportable	bool	`json:"exportable,omitempty"`
	// AllowPlaintextBackup - If set, enables taking backup of named key in the plaintext format. Once set, this cannot
	// be disabled.
	AllowPlaintextBackup	bool	`json:"allow_plaintext_backup,omitempty"`
}

// GitHubAuthResponse is a response for github auth.
type GitHubAuthResponse struct {
	LeaseID		string			`json:"lease_id,omitempty"`
	Renewable	bool			`json:"renewable,omitempty"`
	LeaseDuration	int64			`json:"lease_duration,omitempty"`
	Data		map[string]interface{}	`json:"data,omitempty"`
	Warnings	map[string]interface{}	`json:"warnings,omitempty"`
	Auth		struct {
		ClientToken	string		`json:"client_token,omitempty"`
		Accessor	string		`json:"accessor,omitempty"`
		Policies	[]string	`json:"policies,omitempty"`
		Metadata	struct {
			Username	string	`json:"username,omitempty"`
			Org		string	`json:"org,omitempty"`
		}	`json:"metadata"`
	}	`json:"auth"`
}

// AWSAuthResponse is a response for github auth.
type AWSAuthResponse struct {
	LeaseID		string			`json:"lease_id,omitempty"`
	Renewable	bool			`json:"renewable,omitempty"`
	LeaseDuration	int64			`json:"lease_duration,omitempty"`
	Data		map[string]interface{}	`json:"data,omitempty"`
	Warnings	map[string]interface{}	`json:"warnings,omitempty"`
	Auth		struct {
		ClientToken	string		`json:"client_token,omitempty"`
		Accessor	string		`json:"accessor,omitempty"`
		Policies	[]string	`json:"policies,omitempty"`
		Metadata	struct {
			RoleTagMaxTTL	string	`json:"role_tag_max_ttl,omitempty"`
			InstanceID	string	`json:"instance_id,omitempty"`
			AMIID		string	`json:"ami_id,omitempty"`
			Role		string	`json:"role,omitempty"`
			AuthType	string	`json:"auth_type,omitempty"`
		}	`json:"metadata"`
	}	`json:"auth"`
	Errors	[]string	`json:"errors,omitempty"`
}
