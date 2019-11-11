package secrets

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMockClient(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()

	err := client.Put("testo", map[string]string{"key_123": "value_xyz"})
	assert.Nil(err)

	vals, err := client.Get("testo")
	assert.Nil(err)
	assert.Equal("value_xyz", vals["key_123"])

	_, err = client.Get("fake_test")
	assert.NotNil(err)

	err = client.Delete("another_fake")
	assert.NotNil(err)

	err = client.Delete("testo")
	assert.Nil(err)

	_, err = client.Get("testo")
	assert.NotNil(err)
}

func TestMockClientTransitEncrypt(t *testing.T) {
	assert := assert.New(t)
	client := NewMockClient()

	err := client.CreateTransitKey(context.TODO(), "key1", map[string]interface{}{"mock_option": true})
	assert.Nil(err)

	cipher, err := client.Encrypt(context.TODO(), "key1", []byte(""), []byte("testo"))
	assert.Nil(err)
	assert.NotEmpty(string(cipher))

	// Decrypt with correct context
	plaintext, err := client.Decrypt(context.TODO(), "key1", []byte(""), cipher)
	assert.Nil(err)
	assert.Equal("testo", plaintext)

	// Decrypt with incorrect context
	plaintext, err = client.Decrypt(context.TODO(), "key1", []byte("bad"), cipher)
	assert.Nil(err)
	assert.NotEqual("testo", plaintext)
}

func TestMockClientTransitKeyOperations(t *testing.T) {
	assert := assert.New(t)
	client := NewMockClient()

	err := client.CreateTransitKey(context.TODO(), "key1", map[string]interface{}{"mock_option": true})
	assert.Nil(err)

	// Error when deleting a non deletion_allowed key
	err = client.DeleteTransitKey(context.TODO(), "key1")
	assert.NotNil(err)

	// Configure Key
	err = client.ConfigureTransitKey(context.TODO(), "key1", map[string]interface{}{"deletion_allowed": true})
	assert.Nil(err)

	// Read Key
	keyData, err := client.ReadTransitKey(context.TODO(), "key1")
	assert.Nil(err)
	assert.Equal(true, keyData["deletion_allowed"])

	// Successfully delete key
	err = client.DeleteTransitKey(context.TODO(), "key1")
	assert.Nil(err)
}

func TestMockClientTransitNoKeyFailures(t *testing.T) {
	assert := assert.New(t)
	client := NewMockClient()

	// Error when deleting a nonexistent key
	err := client.DeleteTransitKey(context.TODO(), "key1")
	assert.NotNil(err)

	// Error configuring nonexistent key
	err = client.ConfigureTransitKey(context.TODO(), "key1", map[string]interface{}{"deletion_allowed": true})
	assert.NotNil(err)

	// Error reading nonexistent key
	_, err = client.ReadTransitKey(context.TODO(), "key1")
	assert.NotNil(err)

	// Error encrypting using nonexistent key
	_, err = client.Encrypt(context.TODO(), "key1", []byte(""), []byte("testo"))
	assert.NotNil(err)
}
