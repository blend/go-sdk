package testutil

import (
	"context"
	"encoding/base64"

	"github.com/blend/go-sdk/vault"
)

// Assert MockTransitClient implements TransitClient
var (
	_ vault.Client        = MockTransitClient{}
	_ vault.TransitClient = MockTransitClient{}
)

// MockTransitClient skips interactions with the vault for encryption/decryption
type MockTransitClient struct {
	vault.Client
}

// Encrypt just returns the input in the mock
func (m MockTransitClient) Encrypt(ctx context.Context, key string, context, data []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(data), nil
}

// Decrypt just returns the input in the mock
func (m MockTransitClient) Decrypt(ctx context.Context, key string, context []byte, ciphertext string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(ciphertext)
}
