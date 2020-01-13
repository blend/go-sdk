package secrets

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
