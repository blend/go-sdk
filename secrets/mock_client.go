package secrets

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/util"
)

var _ Client = &MockClient{}
var _ TransitClient = &MockClient{}

// NewMockClient creates a new mock client.
func NewMockClient() *MockClient {
	return &MockClient{
		SecretValues: make(map[string]Values),
		TransitKeys:  make(map[string]map[string][]byte),
	}
}

// MockClient is a mock events client
type MockClient struct {
	SecretValues map[string]Values
	TransitKeys  map[string]map[string][]byte
}

// Put puts a value.
func (c *MockClient) Put(key string, data Values, options ...Option) error {
	c.SecretValues[key] = data

	return nil
}

// Get gets a value at a given key.
func (c *MockClient) Get(key string, options ...Option) (Values, error) {
	val, exists := c.SecretValues[key]
	if !exists {
		return nil, fmt.Errorf("Key not found: %s", key)
	}

	return val, nil
}

// Delete deletes a key.
func (c *MockClient) Delete(key string, options ...Option) error {
	if _, exists := c.SecretValues[key]; !exists {
		return fmt.Errorf("Key not found: %s", key)
	}

	delete(c.SecretValues, key)

	return nil
}

// CreateTransitKey creates a new transit key.
func (c *MockClient) CreateTransitKey(name string) {
	c.TransitKeys[name] = make(map[string][]byte)
}

func (c *MockClient) deriveTransitKey(name string, context []byte) ([]byte, error) {
	contextStr := string(context)

	keyPath, ok := c.TransitKeys[name]
	if !ok {
		return nil, fmt.Errorf("No key")
	}

	key, ok := keyPath[contextStr]
	if !ok {
		key, _ = util.Crypto.CreateKey(32)
		c.TransitKeys[name][contextStr] = key
	}

	return key, nil
}

// Encrypt encrypts a given set of data.
func (c *MockClient) Encrypt(ctx context.Context, name string, context, data []byte) (string, error) {
	key, err := c.deriveTransitKey(name, context)
	if err != nil {
		return "", err
	}

	encryptedData, err := util.Crypto.Encrypt(key, data)
	return string(encryptedData), err
}

// Decrypt decrypts a given set of data.
func (c *MockClient) Decrypt(ctx context.Context, name string, context []byte, ciphertext string) ([]byte, error) {
	key, err := c.deriveTransitKey(name, context)
	if err != nil {
		return nil, err
	}

	return util.Crypto.Decrypt(key, []byte(ciphertext))
}
