package cache

import (
	"context"
	"reflect"
	"sync"
	"time"
	"unsafe"
	"github.com/blend/go-sdk/async"
)

var (
	_ Cache  = (*LocalCache)(nil)
	_ Locker = (*LocalCache)(nil)
)

// NewLocalCache returns a new LocalLocalCache.
func NewLocalCache(options ...LocalCacheOption) *LocalCache {
	c := LocalCache{
		Data: make(map[interface{}]*Value),
		LRU:  NewLRUQueue(),
	}
	c.Sweeper = async.NewInterval(c.Sweep, 500*time.Millisecond)
	for _, opt := range options {
		opt(&c)
	}
	return &c
}

// LocalCacheOption is a local cache option.
type LocalCacheOption func(*LocalCache)

// OptLocalCacheSweepInterval sets the local cache sweep interval.
func OptLocalCacheSweepInterval(d time.Duration) LocalCacheOption {
	return func(lc *LocalCache) {
		lc.Sweeper = async.NewInterval(lc.Sweep, d)
	}
}

// LocalCache is a memory LocalCache.
type LocalCache struct {
	sync.RWMutex
	Data    map[interface{}]*Value
	LRU     LRU
	Sweeper *async.Interval
}

// Start starts the sweeper.
func (lc *LocalCache) Start() error {
	return lc.Sweeper.Start()
}

// NotifyStarted returns the underlying started signal.
func (lc *LocalCache) NotifyStarted() <-chan struct{} {
	return lc.Sweeper.NotifyStarted()
}

// Stop stops the sweeper.
func (lc *LocalCache) Stop() error {
	return lc.Sweeper.Stop()
}

// NotifyStopped returns the underlying stopped signal.
func (lc *LocalCache) NotifyStopped() <-chan struct{} {
	return lc.Sweeper.NotifyStopped()
}

type removeHandler struct {
	Key     interface{}
	Handler func(interface{}, RemovalReason)
}

// Sweep checks keys for expired ttls.
// If any values are configured with 'OnSweep' handlers, they will be called
// outside holding the critical section.
func (lc *LocalCache) Sweep(ctx context.Context) error {
	lc.Lock()
	now := time.Now().UTC()

	var keysToRemove []interface{}
	var handlers []removeHandler
	lc.LRU.Consume(func(v *Value) bool {
		if !v.Expires.IsZero() && now.After(v.Expires) {
			keysToRemove = append(keysToRemove, v.Key)
			if v.OnRemove != nil {
				handlers = append(handlers, removeHandler{
					Key:     v.Key,
					Handler: v.OnRemove,
				})
			}
			return true
		}
		return false
	})

	for _, key := range keysToRemove {
		delete(lc.Data, key)
	}
	lc.Unlock()

	// call the handlers outside the critical section.
	for _, handler := range handlers {
		handler.Handler(handler.Key, Expired)
	}
	return nil
}

// Set adds a LocalCache item.
func (lc *LocalCache) Set(key, value interface{}, options ...ValueOption) {
	if key == nil {
		panic("local cache: nil key")
	}

	if !reflect.TypeOf(key).Comparable() {
		panic("local cache: key is not comparable")
	}

	v := Value{
		Timestamp: time.Now().UTC(),
		Key:       key,
		Value:     value,
	}

	for _, opt := range options {
		opt(&v)
	}

	lc.Lock()
	if lc.Data == nil {
		lc.Data = make(map[interface{}]*Value)
	}
	if value, ok := lc.Data[key]; ok {
		lc.LRU.Fix(&v)
		*value = v
	} else {
		lc.Data[key] = &v
		lc.LRU.Push(&v)
	}
	lc.Unlock()
}

// Get gets a value based on a key.
func (lc *LocalCache) Get(key interface{}) (value interface{}, hit bool) {
	lc.RLock()
	valueNode, ok := lc.Data[key]
	lc.RUnlock()
	if ok {
		value = valueNode.Value
		hit = true
		return
	}
	return
}

// GetOrSet gets a value by a key, and in the case of a miss, sets the value from a given value provider lazily.
// Hit indicates that the provider was not called.
func (lc *LocalCache) GetOrSet(key interface{}, valueProvider func() (interface{}, error), options ...ValueOption) (value interface{}, hit bool, err error) {
	if key == nil {
		panic("local cache: nil key")
	}

	if !reflect.TypeOf(key).Comparable() {
		panic("local cache: key is not comparable")
	}

	// check if we already have the value
	lc.RLock()
	valueNode, ok := lc.Data[key]
	lc.RUnlock()

	if ok {
		value = valueNode.Value
		hit = true
		return
	}

	// call the value provider outside the critical section.
	// this will create a meaningful gap between releasing the
	// read lock and grabbing the write lock.
	value, err = valueProvider()
	if err != nil {
		return
	}

	// we didn't have the value, grab the write lock
	lc.Lock()
	defer lc.Unlock()

	// double checked locks for the children
	// we do this because there may have been a write while we waited
	// for the exclusive lock.
	valueNode, ok = lc.Data[key]
	if ok {
		value = valueNode.Value
		hit = true
		return
	}

	// set up the value
	v := Value{
		Timestamp: time.Now().UTC(),
		Key:       key,
		Value:     value,
	}
	// apply options
	for _, opt := range options {
		opt(&v)
	}

	// upsert
	if value, ok := lc.Data[key]; ok {
		lc.LRU.Fix(&v)
		*value = v
	} else {
		lc.Data[key] = &v
		lc.LRU.Push(&v)
	}

	return
}

// Has returns if the key is present in the LocalCache.
func (lc *LocalCache) Has(key interface{}) (has bool) {
	lc.RLock()
	_, has = lc.Data[key]
	lc.RUnlock()
	return
}

// Remove removes a specific key.
func (lc *LocalCache) Remove(key interface{}) (value interface{}, hit bool) {
	lc.Lock()
	valueData, ok := lc.Data[key]
	if ok {
		delete(lc.Data, key)
		lc.LRU.Remove(key)
	}
	lc.Unlock()
	if !ok {
		return
	}

	value = valueData.Value
	hit = true
	if valueData.OnRemove != nil {
		valueData.OnRemove(key, Removed)
	}
	return
}

// Stats returns the LocalCache stats.
func (lc *LocalCache) Stats() (stats Stats) {
	lc.RLock()
	defer lc.RUnlock()

	stats.Count = len(lc.Data)
	now := time.Now().UTC()
	for _, item := range lc.Data {
		age := now.Sub(item.Timestamp)
		if stats.MaxAge < age {
			stats.MaxAge = age
		}
		stats.SizeBytes += int(unsafe.Sizeof(item))
	}
	return
}
