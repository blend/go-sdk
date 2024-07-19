/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import "time"

const (
	// DefaultAddr is the default addr.
	DefaultAddr = "http://127.0.0.1:8200"
	// DefaultTimeout is the default timeout.
	DefaultTimeout = time.Second
	// DefaultMount is the default kv mount.
	DefaultMount = "/secret"
)

const (
	// EnvVarVaultAddr is the environment variable for the vault address.
	EnvVarVaultAddr = "VAULT_ADDR"
	// EnvVarVaultMount is the environment variable for the vault mount.
	EnvVarVaultMount = "VAULT_MOUNT"
	// EnvVarVaultToken is the environment variable for the vault token.
	EnvVarVaultToken = "VAULT_TOKEN"
	// EnvVarVaultCertAuthorityPath is the environment variable for the vault certificate authority.
	EnvVarVaultCertAuthorityPath = "VAULT_CACERT"
	// EnvVarVaultTimeout is the environment variable for how long to wait for vault to timeout. The values here
	// are parsed by time.ParseDuration. Examples (5s = five seconds, 100ms = 100 milliseconds, etc.)
	EnvVarVaultTimeout = "VAULT_TIMEOUT"
)

const (
	// MethodGet is a request method.
	MethodGet = "GET"
	// MethodPost is a request method.
	MethodPost = "POST"
	// MethodPut is a request method.
	MethodPut = "PUT"
	// MethodDelete is a request method.
	MethodDelete = "DELETE"
	// MethodList is a request method.
	MethodList = "LIST"

	// HeaderVaultToken is the vault token header.
	HeaderVaultToken = "X-Vault-Token"
	// HeaderContentType is the content type header.
	HeaderContentType = "Content-Type"
	// ContentTypeApplicationJSON is a content type.
	ContentTypeApplicationJSON = "application/json"

	// DefaultBufferPoolSize is the default buffer pool size.
	DefaultBufferPoolSize = 1024

	// ReflectTagName is a reflect tag name.
	ReflectTagName = "secret"

	// Version1 is a constant.
	Version1 = "1"
	// Version2 is a constant.
	Version2 = "2"
)

// These types are encryption algorithms that can be used when creating a transit key
const (
	TypeAES256GCM96      = "aes256-gcm96"
	TypeCHACHA20POLY1305 = "chacha20-poly1305"
	TypeED25519          = "ed25519"
	TypeECDSAP256        = "ecdsa-p256"
	TypeRSA2048          = "rsa-2048"
	TypeRSA4096          = "rsa-4096"
)

// These constants are used to sign the get identity request
const (
	// STSURL is the url of the sts call
	STSURL = "https://sts.amazonaws.com"
	// STSGetIdentityBody is the body of the post request
	STSGetIdentityBody = "Action=GetCallerIdentity&Version=2011-06-15"
)

// constants required for login /v1/auth/aws/login
const (
	// AWSAuthLoginPath is the login path for aws iam auth
	AWSAuthLoginPath = "/v1/auth/aws/login"
)
