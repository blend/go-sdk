package crypto

import (
	"sync"
	"time"
)

// KeyCache is a cache for encryption keys.
type KeyCache struct {
	sync.Mutex
	KeyProvider func(string) ([]byte, error)
	Keys        map[string]CachedKey
}

// GetKey gets a cached key or fetches a new one.
func (kc *KeyCache) GetKey(context string) ([]byte, error) {
	kc.Lock()
	defer kc.Unlock()

	if kc.Keys == nil {
		kc.Keys = make(map[string]CachedKey)
	}

	if key, ok := kc.Keys[context]; ok {
		return key.Key, nil
	}

	key, err := kc.KeyProvider(context)
	if err != nil {
		return nil, err
	}

	kc.Keys[context] = CachedKey{
		Context: context,
		Added:   time.Now().UTC(),
		Key:     key,
	}
	return key, nil
}

// PurgeKeys purges keys added before a given time.
func (kc *KeyCache) PurgeKeys(addedBefore time.Time) {
	kc.Lock()
	defer kc.Unlock()

	if kc.Keys == nil {
		kc.Keys = make(map[string]CachedKey)
	}

	var purged []string
	for context, key := range kc.Keys {
		if key.Added.Before(addedBefore) {
			purged = append(purged, context)
		}
	}

	for _, context := range purged {
		delete(kc.Keys, context)
	}
}

// CachedKey is a versioned key.
type CachedKey struct {
	Context string
	Added   time.Time
	Key     []byte
}
