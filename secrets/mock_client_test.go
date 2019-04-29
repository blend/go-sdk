package secrets

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMockClient(t *testing.T) {
	assert := assert.New(t)
	todo := context.TODO()

	client := NewMockClient()

	err := client.Put(todo, "testo", map[string]string{"key_123": "value_xyz"})
	assert.Nil(err)

	vals, err := client.Get(todo, "testo")
	assert.Nil(err)
	assert.Equal("value_xyz", vals["key_123"])

	_, err = client.Get(todo, "fake_test")
	assert.NotNil(err)

	err = client.Delete(todo, "another_fake")
	assert.NotNil(err)

	err = client.Delete(todo, "testo")
	assert.Nil(err)

	_, err = client.Get(todo, "testo")
	assert.NotNil(err)
}

func TestMockClientList(t *testing.T) {
	assert := assert.New(t)
	todo := context.TODO()
	client := NewMockClient()
	data := map[string]string{"key_123": "value_xyz"}

	f := func(path string) []string {
		vals, _ := client.List(todo, path)
		return vals
	}

	assert.Nil(client.Put(todo, "secret/service/abc/key1", data))
	assert.Nil(client.Put(todo, "secret/service/abc/key2", data))
	assert.Nil(client.Put(todo, "secret/service/abc/key3", data))
	assert.Nil(client.Put(todo, "secret/service/abc/folder1/f1key1", data))
	assert.Nil(client.Put(todo, "secret/service/abc/folder1/f1key2", data))
	assert.Nil(client.Put(todo, "secret/service/abc/folder2/f2key1", data))
	assert.Nil(client.Put(todo, "secret/service/head", data))

	results := f("secret")
	assert.Len(results, 1)
	assert.True(validate(results, "service/"))

	results = f("secret/")
	assert.Len(results, 1)
	assert.True(validate(results, "service/"))

	results = f("secret/service/")
	assert.Len(results, 2)
	assert.True(validate(results, "abc/", "head"))

	results = f("secret/service/abc")
	assert.Len(results, 5)
	assert.True(validate(results, "key1", "key2", "key3", "folder1/", "folder2/"))
}

func validate(keys []string, values ...string) bool {
	m := make(map[string]struct{})

	for _, k := range keys {
		if _, ok := m[k]; ok {
			// keys should never contain duplicates
			return false
		}
		m[k] = struct{}{}
	}

	for _, v := range values {
		if _, ok := m[v]; !ok {
			// every value should be in the set
			return false
		}
	}
	return true
}

func TestMockClientTransit(t *testing.T) {
	assert := assert.New(t)
	client := NewMockClient()

	client.CreateTransitKey(context.TODO(), "key1", map[string]interface{}{})

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

	// Delete Key
	client.DeleteTransitKey(context.TODO(), "key1")

	// Get an error while encrypting with deleted key
	cipher, err = client.Encrypt(context.TODO(), "key1", []byte(""), []byte("testo"))
	assert.NotNil(err)
}
