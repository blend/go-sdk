package secrets

import (
	"context"
	"fmt"
	"strings"

	"github.com/blend/go-sdk/crypto"
)

var (
	_ Client = (*MockClient)(nil)
)

// NewMockClient creates a new mock client.
func NewMockClient() *MockClient {
	return &MockClient{
		SecretValues: make(map[string]Values),
		TransitKeys:  make(map[string][]byte),
	}
}

// MockClient is a mock events client
type MockClient struct {
	SecretValues map[string]Values
	TransitKeys  map[string][]byte
}

// Put puts a value.
func (c *MockClient) Put(_ context.Context, key string, data Values, options ...RequestOption) error {
	c.SecretValues[key] = data

	return nil
}

// Get gets a value at a given key.
func (c *MockClient) Get(_ context.Context, key string, options ...RequestOption) (Values, error) {
	val, exists := c.SecretValues[key]
	if !exists {
		return nil, fmt.Errorf("Key not found: %s", key)
	}

	return val, nil
}

// Delete deletes a key.
func (c *MockClient) Delete(_ context.Context, key string, options ...RequestOption) error {
	if _, exists := c.SecretValues[key]; !exists {
		return fmt.Errorf("Key not found: %s", key)
	}

	delete(c.SecretValues, key)
	return nil
}

// List lists keys on a path
func (c *MockClient) List(_ context.Context, path string, options ...RequestOption) ([]string, error) {
	keys := make([]string, 0)
	folderSet := make(map[string]struct{})
	p := path
	if !strings.HasSuffix(path, "/") {
		p = path + "/"
	}
	for k := range c.SecretValues {
		if strings.HasPrefix(k, p) {
			s := strings.TrimPrefix(k, p)
			if strings.ContainsRune(s, '/') {
				folder := fmt.Sprintf("%s/", strings.Split(s, "/")[0])
				if _, ok := folderSet[folder]; !ok {
					folderSet[folder] = struct{}{}
					keys = append(keys, folder)
				}
			} else {
				keys = append(keys, s)
			}
		}
	}
	return keys, nil
}

// CreateTransitKey creates a new transit key.
func (c *MockClient) CreateTransitKey(name string) error {
	key, err := crypto.CreateKey(32)
	c.TransitKeys[name] = key

	return err
}

// Encrypt encrypts a given set of data.
func (c *MockClient) Encrypt(ctx context.Context, name string, context, data []byte) (string, error) {
	_, ok := c.TransitKeys[name]
	if !ok {
		return "", fmt.Errorf("No key")
	}

	encryptedData, err := crypto.Encrypt(c.TransitKeys[name], data)
	return string(encryptedData), err
}

// Decrypt decrypts a given set of data.
func (c *MockClient) Decrypt(ctx context.Context, name string, context []byte, ciphertext string) ([]byte, error) {
	_, ok := c.TransitKeys[name]
	if !ok {
		return nil, fmt.Errorf("No key")
	}

	return crypto.Decrypt(c.TransitKeys[name], []byte(ciphertext))
}
