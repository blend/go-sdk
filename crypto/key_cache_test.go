package crypto

import (
	"fmt"
	"testing"
	"time"
	"github.com/blend/go-sdk/assert"
)

func TestKeyCache(t *testing.T) {
	assert := assert.New(t)

	cache := &KeyCache{
		KeyProvider: func(context string) ([]byte, error) {
			return []byte(context), nil
		},
	}

	found, err := cache.GetKey("foo")
	assert.Nil(err)
	assert.Equal([]byte("foo"), found)
	assert.Len(cache.Keys, 1)

	found, err = cache.GetKey("bar")
	assert.Nil(err)
	assert.Equal([]byte("bar"), found)
	assert.Len(cache.Keys, 2)
}

func TestKeyCacheRespectsErrors(t *testing.T) {
	assert := assert.New(t)

	cache := &KeyCache{
		KeyProvider: func(context string) ([]byte, error) {
			return nil, fmt.Errorf("this is only a test")
		},
	}

	found, err := cache.GetKey("foo")
	assert.NotNil(err)
	assert.Nil(found)
	assert.Empty(cache.Keys)

	found, err = cache.GetKey("bar")
	assert.NotNil(err)
	assert.Nil(found)
	assert.Empty(cache.Keys)
}

func TestKeyCachePurgeKeys(t *testing.T) {
	assert := assert.New(t)

	cache := &KeyCache{
		KeyProvider: func(context string) ([]byte, error) {
			return []byte(context), nil
		},
		Keys: map[string]CachedKey{
			"foo": CachedKey{Added: time.Date(2019, 04, 10, 01, 02, 03, 04, time.UTC)},
			"bar": CachedKey{Added: time.Date(2019, 04, 11, 01, 02, 03, 04, time.UTC)},
			"baz": CachedKey{Added: time.Date(2019, 04, 12, 01, 02, 03, 04, time.UTC)},
		},
	}
	assert.Len(cache.Keys, 3)
	cache.PurgeKeys(time.Date(2019, 04, 11, 01, 02, 03, 03, time.UTC))
	assert.Len(cache.Keys, 2)

	_, ok := cache.Keys["bar"]
	assert.True(ok)

	_, ok = cache.Keys["baz"]
	assert.True(ok)
}
